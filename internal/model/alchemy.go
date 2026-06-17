package model

// Transaction = eth_getTransactionByHash
type Transaction struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Input    string `json:"input"`
	Value    string `json:"value"`
	GasPrice string `json:"gasPrice"`
	GasLimit string `json:"gas"`
	Nonce    string `json:"nonce"`
}

// RawLog - Ethereum log (receipt.logs)
type RawLog struct {
	Address     string   `json:"address"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
	LogIndex    string   `json:"logIndex"`
	TxHash      string   `json:"transactionHash"`
	BlockNumber string   `json:"blockNumber"`
	BlockHash   string   `json:"blockHash"`
	Removed     bool     `json:"removed"`
}

// Receipt = eth_getTransactionReceipt
type Receipt struct {
	TxHash          string   `json:"transactionHash"`
	Status          string   `json:"status"`
	GasUsed         string   `json:"gasUsed"`
	ContractAddress string   `json:"contractAddress"`
	Logs            []RawLog `json:"logs"`
}
