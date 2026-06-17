package collector

import (
	"math/big"
	"realtime_blockchain/internal/model"
	"strings"
	"sync"
)

type TokenMetaProvider interface {
	Get(address string) (*model.TokenMeta, bool)
	Upsert(meta *model.TokenMeta)
}

type PriceProvider interface {
	GetUSD(chain string, symbol string, address string) (float64, bool)
}

type Enricher struct {
	in            <-chan model.BlockchainEvent
	out           chan model.EnrichedEvent
	metaProvider  TokenMetaProvider
	priceProvider PriceProvider
	cache         map[string]*model.TokenMeta
	mu            sync.RWMutex
}

func NewEnricher(
	in <-chan model.BlockchainEvent,
	metaProvider TokenMetaProvider,
	priceProvider PriceProvider,
) *Enricher {

	return &Enricher{
		in:            in,
		out:           make(chan model.EnrichedEvent, 1000),
		metaProvider:  metaProvider,
		priceProvider: priceProvider,
		cache:         make(map[string]*model.TokenMeta),
	}
}

func (e *Enricher) Out() <-chan model.EnrichedEvent {
	return e.out
}

func (e *Enricher) Run() {
	defer close(e.out)

	for ev := range e.in {

		meta := e.getTokenMeta(ev.TokenAddress)

		amountHuman := e.toHumanAmount(ev.RawAmount, meta.Decimals)

		direction := e.getDirection(ev.TxFrom, ev.From, ev.To)
		eventType := e.getType(ev)

		valueUSD := 0.0
		price, ok := e.priceProvider.GetUSD("ethereum", meta.Symbol, ev.TokenAddress)
		if ok {
			valueUSD = amountHuman * price
		}

		enriched := model.EnrichedEvent{
			TxHash:       ev.TxHash,
			TxFrom:       ev.TxFrom,
			TxTo:         ev.TxTo,
			TxValue:      ev.TxValue,
			GasUsed:      ev.GasUsed,
			BlockNumber:  ev.BlockNumber,
			BlockHash:    ev.BlockHash,
			TokenAddress: ev.TokenAddress,
			Symbol:       meta.Symbol,
			Name:         meta.Name,
			AmountRaw:    ev.RawAmount,
			AmountHuman:  amountHuman,
			From:         ev.From,
			To:           ev.To,
			Type:         eventType,
			Direction:    direction,
			ValueUSD:     valueUSD,
		}

		e.out <- enriched
	}
}

func (e *Enricher) getTokenMeta(addr string) *model.TokenMeta {
	addr = strings.ToLower(addr)

	e.mu.RLock()
	if m, ok := e.cache[addr]; ok {
		e.mu.RUnlock()
		return m
	}
	e.mu.RUnlock()

	meta, ok := e.metaProvider.Get(addr)
	if !ok {
		meta = &model.TokenMeta{
			Address:  addr,
			Symbol:   "UNKNOWN",
			Name:     "UNKNOWN",
			Decimals: 18,
		}
	}

	e.mu.Lock()
	e.cache[addr] = meta
	e.mu.Unlock()

	return meta
}

func (e *Enricher) toHumanAmount(raw string, decimals int) float64 {
	raw = strings.TrimPrefix(raw, "0x")

	bi := new(big.Int)
	bi.SetString(raw, 16)

	if decimals == 0 {
		return float64(bi.Int64())
	}

	divisor := new(big.Float).SetFloat64(float64Pow(10, decimals))
	val := new(big.Float).SetInt(bi)

	result, _ := new(big.Float).Quo(val, divisor).Float64()
	return result
}

func float64Pow(base float64, exp int) float64 {
	result := 1.0
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

func (e *Enricher) getDirection(txFrom, from, to string) model.Direction {
	txFrom = strings.ToLower(txFrom)
	from = strings.ToLower(from)
	to = strings.ToLower(to)

	if txFrom == "" {
		return model.DirectionUnknown
	}

	// contract → user (incoming)
	if to == txFrom {
		return model.DirectionIn
	}

	// user → contract (outgoing)
	if from == txFrom {
		return model.DirectionOut
	}

	// fallback
	return model.DirectionUnknown
}

func (e *Enricher) getType(ev model.BlockchainEvent) model.EventType {
	// Only ERC20 transfer already known
	return model.EventTransfer
}
