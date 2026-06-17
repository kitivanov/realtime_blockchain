package main

import (
	"log"
	"time"

	"realtime_blockchain/internal/collector"
	"realtime_blockchain/internal/config"
	"realtime_blockchain/internal/registry"
	"realtime_blockchain/internal/transport/alchemy"
	"realtime_blockchain/internal/transport/coingecko"
)

func main() {
	cfg := config.Load()

	rpc := alchemy.NewRPCClient(cfg.AlchemyURLHTTP + cfg.AlchemyKey)

	tokenRegistry, tErr := registry.NewTokenRegistryFromFile("data/tokens.json")
	if tErr != nil {
		log.Fatal(tErr)
	}

	priceCache := registry.NewCache()
	coinGecko := coingecko.NewCoinGeckoClient(cfg, priceCache)
	go coinGecko.Start()
	// coinGecko.Start()
	// coinGecko.GetUSD("ethereum", "0x6b175474e89094c44da98b954eedeac495271d0f")

	ws := alchemy.NewWSLogCollector(cfg)
	err := ws.Connect()

	if err != nil {
		log.Fatal(err)
	}

	deduplicator := collector.NewDeduplicator(ws.Out(), 1*time.Minute)
	deduplicator.StartCleanup()
	go deduplicator.Run()

	batcher := collector.NewBatcher(deduplicator.Out(), 50, 200*time.Millisecond)
	go batcher.Run()

	fetcherFanOut := collector.NewFetcherFanOut(batcher.Out())
	go fetcherFanOut.Run()

	receiptFetcher := collector.NewReceiptFetcher(fetcherFanOut.ReceiptOut(), rpc)
	txFetcher := collector.NewTransactionFetcher(fetcherFanOut.TxOut(), rpc)

	go receiptFetcher.Run()
	go txFetcher.Run()

	joiner := collector.NewJoiner(txFetcher.Out(), receiptFetcher.Out())
	go joiner.Run()

	ERC20Decoder := collector.NewERC20Decoder(joiner.Out())
	go ERC20Decoder.Run()

	enricher := collector.NewEnricher(ERC20Decoder.Out(), tokenRegistry, coinGecko)
	go enricher.Run()

	go func() {
		for i := range enricher.Out() {
			collector.HandleEvents(i)
		}
	}()

	err = ws.SubscribeLogs()
	if err != nil {
		log.Fatal(err)
	}

	ws.ListenLogs()
}
