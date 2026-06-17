package alchemy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"realtime_blockchain/internal/model"
)

type rpcReq struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type rpcResp[T any] struct {
	ID     int `json:"id"`
	Result *T  `json:"result"`
	Error  *struct {
		Code    int
		Message string
	} `json:"error,omitempty"`
}

type BatchItem struct {
	ID     int
	TxHash string
}

type ReceiptBatchResult struct {
	ID      int
	Receipt *model.Receipt
	Error   error
}

type TransactionBatchResult struct {
	ID     int
	TxHash string
	Tx     *model.Transaction
	Error  error
}

type RPCClient struct {
	url string
}

func NewRPCClient(url string) *RPCClient {
	return &RPCClient{url: url}
}

func doBatch[T any](url string, reqs []rpcReq) ([]rpcResp[T], error) {

	body, err := json.Marshal(reqs)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []rpcResp[T]

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *RPCClient) GetTxReceiptsBatch(items []BatchItem) ([]ReceiptBatchResult, error) {

	reqs := make([]rpcReq, len(items))

	for i, it := range items {
		reqs[i] = rpcReq{
			JSONRPC: "2.0",
			ID:      it.ID,
			Method:  "eth_getTransactionReceipt",
			Params:  []interface{}{it.TxHash},
		}
	}

	raw, err := doBatch[model.Receipt](r.url, reqs)
	if err != nil {
		return nil, err
	}

	byID := make(map[int]rpcResp[model.Receipt], len(raw))

	for _, r := range raw {
		byID[r.ID] = r
	}

	out := make([]ReceiptBatchResult, len(items))

	for i, it := range items {

		res, ok := byID[it.ID]
		if !ok {
			out[i] = ReceiptBatchResult{
				ID:    it.ID,
				Error: fmt.Errorf("missing response"),
			}
			continue
		}

		if res.Error != nil {
			out[i] = ReceiptBatchResult{
				ID:    it.ID,
				Error: fmt.Errorf("%s", res.Error.Message),
			}
			continue
		}

		if res.Result == nil {
			continue
		}

		out[i] = ReceiptBatchResult{
			ID:      it.ID,
			Receipt: res.Result,
		}
	}

	return out, nil
}

func (r *RPCClient) GetTransactionsBatch(items []BatchItem) ([]TransactionBatchResult, error) {

	reqs := make([]rpcReq, len(items))

	for i, it := range items {
		reqs[i] = rpcReq{
			JSONRPC: "2.0",
			ID:      it.ID,
			Method:  "eth_getTransactionByHash",
			Params:  []interface{}{it.TxHash},
		}
	}

	raw, err := doBatch[model.Transaction](r.url, reqs)
	if err != nil {
		return nil, err
	}

	byID := make(map[int]rpcResp[model.Transaction], len(raw))

	for _, r := range raw {
		byID[r.ID] = r
	}

	out := make([]TransactionBatchResult, len(items))

	for i, it := range items {

		res, ok := byID[it.ID]
		if !ok {
			out[i] = TransactionBatchResult{
				ID:     it.ID,
				TxHash: it.TxHash,
				Error:  fmt.Errorf("missing response"),
			}
			continue
		}

		if res.Error != nil {
			out[i] = TransactionBatchResult{
				ID:     it.ID,
				TxHash: it.TxHash,
				Error:  fmt.Errorf("%s", res.Error.Message),
			}
			continue
		}

		out[i] = TransactionBatchResult{
			ID:     it.ID,
			TxHash: it.TxHash,
			Tx:     res.Result,
		}

	}

	return out, nil
}

// func (r *RPCClient) GetBlockByNumber(number string) (*model.Block, error) {

// 	payload := rpcReq{
// 		JSONRPC: "2.0",
// 		ID:      1,
// 		Method:  "eth_getBlockByNumber",
// 		Params:  []interface{}{number, true},
// 	}

// 	body, _ := json.Marshal(payload)

// 	resp, err := http.Post(r.url, "application/json", bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var wrapper struct {
// 		Result model.Block `json:"result"`
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
// 		return nil, err
// 	}

// 	return &wrapper.Result, nil
// }
