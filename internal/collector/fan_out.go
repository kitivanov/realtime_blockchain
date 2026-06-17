package collector

import (
	"realtime_blockchain/internal/model"
)

type FetcherFanOut struct {
	in <-chan []model.LogTrigger

	receiptOut chan []model.LogTrigger
	txOut      chan []model.LogTrigger
}

func NewFetcherFanOut(in <-chan []model.LogTrigger) *FetcherFanOut {
	return &FetcherFanOut{
		in: in,

		receiptOut: make(chan []model.LogTrigger, 100),
		txOut:      make(chan []model.LogTrigger, 100),
	}
}

func (f *FetcherFanOut) ReceiptOut() <-chan []model.LogTrigger {
	return f.receiptOut
}

func (f *FetcherFanOut) TxOut() <-chan []model.LogTrigger {
	return f.txOut
}

func (f *FetcherFanOut) Run() {
	for batch := range f.in {

		if len(batch) == 0 {
			continue
		}

		receiptBatch := make([]model.LogTrigger, len(batch))
		txBatch := make([]model.LogTrigger, len(batch))

		copy(receiptBatch, batch)
		copy(txBatch, batch)

		// receipt pipeline
		select {
		case f.receiptOut <- receiptBatch:
		default:
			// TODO
		}

		// tx pipeline
		select {
		case f.txOut <- txBatch:
		default:
			// TODO
		}
	}

	close(f.receiptOut)
	close(f.txOut)
}

type DecoderFanOut struct {
	in <-chan model.JoinedTransaction

	erc20Out    chan model.JoinedTransaction
	swapOut     chan model.JoinedTransaction
	nftOut      chan model.JoinedTransaction
	approvalOut chan model.JoinedTransaction
}

func NewDecoderFanOut(in <-chan model.JoinedTransaction) *DecoderFanOut {
	return &DecoderFanOut{
		in: in,

		erc20Out:    make(chan model.JoinedTransaction, 100),
		swapOut:     make(chan model.JoinedTransaction, 100),
		nftOut:      make(chan model.JoinedTransaction, 100),
		approvalOut: make(chan model.JoinedTransaction, 100),
	}
}

func (f *DecoderFanOut) ERC20Out() <-chan model.JoinedTransaction {
	return f.erc20Out
}

func (f *DecoderFanOut) SwapOut() <-chan model.JoinedTransaction {
	return f.swapOut
}

func (f *DecoderFanOut) NFTOut() <-chan model.JoinedTransaction {
	return f.nftOut
}

func (f *DecoderFanOut) ApprovalOut() <-chan model.JoinedTransaction {
	return f.approvalOut
}

func (f *DecoderFanOut) Run() {
	defer func() {
		close(f.erc20Out)
		close(f.swapOut)
		close(f.nftOut)
		close(f.approvalOut)
	}()

	for joined := range f.in {

		select {
		case f.erc20Out <- joined:
		default:
		}

		// select {
		// case f.swapOut <- joined:
		// default:
		// }

		// select {
		// case f.nftOut <- joined:
		// default:
		// }

		// select {
		// case f.approvalOut <- joined:
		// default:
		// }
	}
}
