package collector

import (
	"sync"

	"realtime_blockchain/internal/model"
	"realtime_blockchain/internal/transport/alchemy"
)

type Joiner struct {
	txIn       <-chan []alchemy.TransactionBatchResult
	receiptIn  <-chan []alchemy.ReceiptBatchResult
	out        chan model.JoinedTransaction
	mu         sync.Mutex
	txMap      map[string]*model.Transaction
	receiptMap map[string]*model.Receipt
}

func NewJoiner(
	txIn <-chan []alchemy.TransactionBatchResult,
	receiptIn <-chan []alchemy.ReceiptBatchResult,
) *Joiner {

	return &Joiner{
		txIn:       txIn,
		receiptIn:  receiptIn,
		out:        make(chan model.JoinedTransaction, 1000),
		txMap:      make(map[string]*model.Transaction),
		receiptMap: make(map[string]*model.Receipt),
	}
}

func (j *Joiner) Out() <-chan model.JoinedTransaction {
	return j.out
}

func (j *Joiner) Run() {
	defer close(j.out)

	for {
		select {

		case batch, ok := <-j.txIn:
			if !ok {
				j.txIn = nil
			} else {
				j.handleTransactions(batch)
			}

		case batch, ok := <-j.receiptIn:
			if !ok {
				j.receiptIn = nil
			} else {
				j.handleReceipts(batch)
			}
		}

		if j.txIn == nil && j.receiptIn == nil {
			return
		}
	}
}

func (j *Joiner) handleTransactions(batch []alchemy.TransactionBatchResult) {
	j.mu.Lock()
	defer j.mu.Unlock()

	for _, item := range batch {

		if item.Tx == nil {
			continue
		}

		j.txMap[item.TxHash] = item.Tx

		j.tryJoin(item.TxHash)
	}
}

func (j *Joiner) handleReceipts(batch []alchemy.ReceiptBatchResult) {
	j.mu.Lock()
	defer j.mu.Unlock()

	for _, item := range batch {

		if item.Receipt == nil {
			continue
		}

		txHash := item.Receipt.TxHash

		j.receiptMap[txHash] = item.Receipt

		j.tryJoin(txHash)
	}
}

func (j *Joiner) tryJoin(txHash string) {

	tx, okTx := j.txMap[txHash]
	receipt, okReceipt := j.receiptMap[txHash]

	if !okTx || !okReceipt {
		return
	}

	j.out <- model.JoinedTransaction{
		Tx:      tx,
		Receipt: receipt,
	}

	delete(j.txMap, txHash)
	delete(j.receiptMap, txHash)
}
