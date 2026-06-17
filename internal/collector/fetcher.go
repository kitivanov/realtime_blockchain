package collector

import (
	"log"
	"realtime_blockchain/internal/model"
	"realtime_blockchain/internal/transport/alchemy"
)

type ReceiptFetcher struct {
	in  <-chan []model.LogTrigger
	out chan []alchemy.ReceiptBatchResult
	rpc *alchemy.RPCClient
}

func NewReceiptFetcher(
	in <-chan []model.LogTrigger,
	rpc *alchemy.RPCClient,
) *ReceiptFetcher {
	return &ReceiptFetcher{
		in:  in,
		out: make(chan []alchemy.ReceiptBatchResult, 100),
		rpc: rpc,
	}
}

func (f *ReceiptFetcher) Out() <-chan []alchemy.ReceiptBatchResult {
	return f.out
}

func (f *ReceiptFetcher) Run() {
	for batch := range f.in {
		log.Println("batch size Receipt IN:", len(batch)) // TODO DELETE THIS

		items := make([]alchemy.BatchItem, len(batch))

		for i, tx := range batch {
			items[i] = alchemy.BatchItem{
				ID:     i,
				TxHash: tx.TxHash,
			}
		}

		results, err := f.rpc.GetTxReceiptsBatch(items)
		if err != nil {
			continue
		}

		f.out <- results
		log.Println("batch size Receipt OUT:", len(results)) // TODO DELETE THIS
	}
}

type TransactionFetcher struct {
	in  <-chan []model.LogTrigger
	out chan []alchemy.TransactionBatchResult
	rpc *alchemy.RPCClient
}

func NewTransactionFetcher(
	in <-chan []model.LogTrigger,
	rpc *alchemy.RPCClient,
) *TransactionFetcher {
	return &TransactionFetcher{
		in:  in,
		out: make(chan []alchemy.TransactionBatchResult, 100),
		rpc: rpc,
	}
}

func (f *TransactionFetcher) Out() <-chan []alchemy.TransactionBatchResult {
	return f.out
}

func (f *TransactionFetcher) Run() {
	for batch := range f.in {
		log.Println("batch size Transaction IN:", len(batch)) // TODO DELETE THIS

		items := make([]alchemy.BatchItem, len(batch))

		for i, tx := range batch {
			items[i] = alchemy.BatchItem{
				ID:     i,
				TxHash: tx.TxHash,
			}
		}

		results, err := f.rpc.GetTransactionsBatch(items)
		if err != nil {
			continue
		}

		f.out <- results
		log.Println("batch size Transaction OUT:", len(results)) // TODO DELETE THIS
	}
}
