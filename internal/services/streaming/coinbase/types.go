package coinbase

import (
	"math/big"
	"time"
)

type SubscribeRequest struct {
	Type       string    `json:"type"`
	ProductIds []string  `json:"product_ids"`
	Channels   []Channel `json:"channels"`
}

type Channel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

type Feed struct {
	Type         string     `json:"type"`
	TradeID      int        `json:"trade_id"`
	MakerOrderID string     `json:"maker_order_id"`
	TakerOrderID string     `json:"taker_order_id"`
	Side         string     `json:"side"`
	Size         *big.Float `json:"size"`
	Price        *big.Float `json:"price"`
	ProductID    string     `json:"product_id"`
	Sequence     int64      `json:"sequence"`
	Time         time.Time  `json:"time"`
	Reason       string     `json:"reason,omitempty"`
}
