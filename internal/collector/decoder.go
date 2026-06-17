package collector

import (
	"math/big"
	"strings"

	"realtime_blockchain/internal/model"
)

const transferSignature = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

func isTransfer(log model.RawLog) bool {
	if len(log.Topics) == 0 {
		return false
	}
	return strings.EqualFold(log.Topics[0], transferSignature)
}

func decodeAddress(topic string) string {
	// topic = 32 bytes, address = last 20 bytes
	return "0x" + topic[26:]
}
func decodeAmount(data string) string {
	data = strings.TrimPrefix(data, "0x")
	b, ok := new(big.Int).SetString(data, 16)
	if !ok {
		return "0"
	}
	return b.String()
}

type ERC20Decoder struct {
	in  <-chan model.JoinedTransaction
	out chan model.BlockchainEvent
}

func NewERC20Decoder(in <-chan model.JoinedTransaction) *ERC20Decoder {
	return &ERC20Decoder{
		in:  in,
		out: make(chan model.BlockchainEvent, 100),
	}
}

func (d *ERC20Decoder) Out() <-chan model.BlockchainEvent {
	return d.out
}

func (d *ERC20Decoder) Run() {
	defer close(d.out)

	for joined := range d.in {

		if joined.Receipt == nil || joined.Tx == nil {
			continue
		}

		for _, rawLog := range joined.Receipt.Logs {

			transfer, ok := DecodeERC20(rawLog)
			if !ok {
				continue
			}

			event := model.BlockchainEvent{
				TxHash:       transfer.TxHash,
				TxFrom:       joined.Tx.From,
				TxTo:         joined.Tx.To,
				TxValue:      joined.Tx.Value,
				GasUsed:      joined.Receipt.GasUsed,
				From:         transfer.From,
				To:           transfer.To,
				TokenAddress: transfer.TokenAddress,
				RawAmount:    transfer.AmountRaw,
				LogIndex:     transfer.LogIndex,
				BlockNumber:  rawLog.BlockNumber,
				BlockHash:    rawLog.BlockHash,
				Removed:      rawLog.Removed,
			}
			d.out <- event
		}
	}
}

func DecodeERC20(log model.RawLog) (*model.TokenTransfer, bool) {

	if len(log.Topics) < 3 {
		return nil, false
	}

	if !strings.EqualFold(log.Topics[0], transferSignature) {
		return nil, false
	}

	return &model.TokenTransfer{
		TxHash:       log.TxHash,
		LogIndex:     log.LogIndex,
		TokenAddress: log.Address,
		From:         decodeAddress(log.Topics[1]),
		To:           decodeAddress(log.Topics[2]),
		AmountRaw:    decodeAmount(log.Data),
	}, true
}
