package model

type TxType string

const (
	TxUnknown       TxType = "unknown"
	TxETHTransfer   TxType = "eth_transfer"
	TxERC20Transfer TxType = "erc20_transfer"
	TxSwap          TxType = "swap"
	TxApproval      TxType = "approval"
)

type EventType string

const (
	EventTransfer EventType = "transfer"
	EventSwap     EventType = "swap"
	EventApproval EventType = "approval"
	EventETH      EventType = "eth_transfer"
)

type Direction string

const (
	DirectionIn      Direction = "in"
	DirectionOut     Direction = "out"
	DirectionSelf    Direction = "self"
	DirectionUnknown Direction = "unknown"
)

type LogTrigger struct {
	TxHash string
}

// TokenTransfer = decoded ERC20 Transfer event
type TokenTransfer struct {
	TxHash       string
	LogIndex     string
	TokenAddress string
	From         string
	To           string
	AmountRaw    string
}

type JoinedTransaction struct {
	Tx      *Transaction
	Receipt *Receipt
}

type BlockchainEvent struct {
	TxHash       string
	TxFrom       string
	TxTo         string
	TxValue      string
	GasUsed      string
	From         string
	To           string
	TokenAddress string
	RawAmount    string
	BlockNumber  string
	BlockHash    string
	LogIndex     string
	Removed      bool
}

type EnrichedEvent struct {
	TxHash       string
	TxFrom       string
	TxTo         string
	TxValue      string
	GasUsed      string
	BlockNumber  string
	BlockHash    string
	TokenAddress string
	Symbol       string
	Name         string
	AmountRaw    string
	AmountHuman  float64
	Type         EventType
	Direction    Direction
	From         string
	To           string
	ValueUSD     float64
}
