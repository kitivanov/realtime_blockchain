package model

import "time"

// raw response from CoinGecko
type CoinGeckoResponse map[string]struct {
	USD float64 `json:"usd"`
}

type PriceEntry struct {
	USD       float64
	UpdatedAt time.Time
	TTL       time.Duration
	NotFound  bool
}

func (p PriceEntry) IsExpired() bool {
	return time.Since(p.UpdatedAt) > p.TTL
}
