package alchemy

import (
	"bytes"
	"encoding/json"
	"log"
	"realtime_blockchain/internal/model"
)

func (c *WSLogCollector) SubscribeLogs() error {

	msg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_subscribe",
		"params": []interface{}{
			"logs",
			map[string]interface{}{
				"topics": []interface{}{
					c.cfg.TransferTopic,
				},
			},
		},
	}

	if err := c.conn.WriteJSON(msg); err != nil {
		return err
	}

	var resp map[string]interface{}
	if err := c.conn.ReadJSON(&resp); err != nil {
		return err
	}

	log.Println("Subscribe created. Response:", resp)

	return nil
}

func (c *WSLogCollector) ListenLogs() {
	log.Println("waiting logs...")
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("ws error:", err)
			return
		}

		c.handleRaw(message)
	}
}

func (c *WSLogCollector) handleRaw(message []byte) {

	if !bytes.Contains(message, []byte("params")) {
		return
	}

	var msg struct {
		Params struct {
			Result struct {
				TransactionHash string `json:"transactionHash"`
			} `json:"result"`
		} `json:"params"`
	}

	if err := json.Unmarshal(message, &msg); err != nil {
		return
	}

	txHash := msg.Params.Result.TransactionHash
	if txHash == "" {
		return
	}

	select {
	case c.out <- model.LogTrigger{TxHash: txHash}:
	default:
		// TODO: drop / metric
	}
}
