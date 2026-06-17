package alchemy

import (
	"fmt"
	"log"
	"realtime_blockchain/internal/config"
	"realtime_blockchain/internal/model"

	"github.com/gorilla/websocket"
)

type WSLogCollector struct {
	conn *websocket.Conn
	cfg  *config.Config
	out  chan model.LogTrigger
}

func NewWSLogCollector(cfg *config.Config) *WSLogCollector {
	return &WSLogCollector{
		cfg: cfg,
		out: make(chan model.LogTrigger, 1000),
	}
}

func (c *WSLogCollector) Out() <-chan model.LogTrigger {
	return c.out
}

func (c *WSLogCollector) Connect() error {
	wsUrl := fmt.Sprintf("%s%s", c.cfg.AlchemyURLWS, c.cfg.AlchemyKey)

	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}

	c.conn = conn
	log.Println("WS connected successfully.")
	return nil
}
