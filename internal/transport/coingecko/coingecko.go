package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"realtime_blockchain/internal/config"
	"realtime_blockchain/internal/model"
	"realtime_blockchain/internal/registry"
	"strings"
	"time"
)

type CoinGeckoClient struct {
	httpClient      *http.Client
	apiURL          string
	cache           *registry.Cache
	apiKey          string
	pages           int
	perPage         int
	ttl             time.Duration
	refreshInterval time.Duration

	stop chan struct{}
}

func NewCoinGeckoClient(cfg *config.Config, cache *registry.Cache) *CoinGeckoClient {
	ttl := time.Duration(cfg.CoinGeckoTTL) * time.Second

	refresh := ttl * 9 / 10

	if refresh <= 0 {
		refresh = ttl / 2
	}

	return &CoinGeckoClient{
		httpClient:      &http.Client{Timeout: 5 * time.Second},
		apiURL:          cfg.CoinGeckoURL,
		cache:           cache,
		apiKey:          cfg.CoinGeckoAPIKey,
		pages:           cfg.CoinGeckoPages,
		perPage:         cfg.CoinGeckoPerPage,
		ttl:             ttl,
		refreshInterval: refresh,
		stop:            make(chan struct{}),
	}
}

func (c *CoinGeckoClient) Start() {
	ticker := time.NewTicker(c.refreshInterval)

	c.refreshMarkets()

	for {
		select {
		case <-ticker.C:
			c.refreshMarkets()

		case <-c.stop:
			ticker.Stop()
			return
		}
	}
}

func (c *CoinGeckoClient) Stop() {
	close(c.stop)
}

func (c *CoinGeckoClient) refreshMarkets() {
	for page := 1; page <= c.pages; page++ {

		url := fmt.Sprintf(
			"%s/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=%d&page=%d&sparkline=false",
			c.apiURL,
			c.perPage,
			page,
		)

		resp, err := c.httpClient.Get(url)
		if err != nil {
			continue
		}

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			log.Println("CoinGecko HTTP error:", resp.StatusCode, string(body))
			continue
		}

		var data []struct {
			ID           string  `json:"id"`
			Symbol       string  `json:"symbol"`
			CurrentPrice float64 `json:"current_price"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			log.Println("CoinGecko decode error:", err)
			continue
		}
		resp.Body.Close()

		log.Println(url)       // todo DELETE THIS
		log.Println(len(data)) // todo DELETE THIS

		now := time.Now()

		// for _, token := range data {
		// 	c.cache.Set(token.ID, model.PriceEntry{
		// 		USD:       token.CurrentPrice,
		// 		UpdatedAt: now,
		// 		TTL:       c.ttl,
		// 		NotFound:  false,
		// 	})
		// }
		for _, token := range data {
			key := strings.ToLower(token.Symbol)

			c.cache.Set(key, model.PriceEntry{
				USD:       token.CurrentPrice,
				UpdatedAt: now,
				TTL:       c.ttl,
				NotFound:  false,
			})
		}
	}
}

func (c *CoinGeckoClient) GetUSD(chain string, symbol string, address string) (float64, bool) {
	address = strings.ToLower(address)
	symbol = strings.ToLower(symbol)

	// 1. cache hit
	if entry, ok := c.cache.Get(symbol); ok {

		// если найден и не истёк
		if !entry.IsExpired() && !entry.NotFound {
			return entry.USD, true
		}

		// если NotFound и TTL не прошёл
		if entry.NotFound && !entry.IsExpired() {
			return 0, false
		}
	}
	return 0, false

	// // 2. fallback direct request
	// price, ok := c.fetchFromAPI(chain, address)
	// if !ok {

	// 	// mark as not found to avoid spam requests
	// 	c.cache.Set(address, model.PriceEntry{ // TODO change key to symbol
	// 		USD:       0,
	// 		UpdatedAt: time.Now(),
	// 		TTL:       c.ttl,
	// 		NotFound:  true,
	// 	})

	// 	return 0, false
	// }

	// // 3. save
	// c.cache.Set(address, model.PriceEntry{ // TODO change key to symbol
	// 	USD:       price,
	// 	UpdatedAt: time.Now(),
	// 	TTL:       c.ttl,
	// 	NotFound:  false,
	// })

	// return price, true
}

func (c *CoinGeckoClient) fetchFromAPI(chain string, address string) (float64, bool) {
	url := fmt.Sprintf(
		"%s/api/v3/simple/token_price/%s?contract_addresses=%s&vs_currencies=usd",
		c.apiURL,
		chain,
		address,
	)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, false
	}
	defer resp.Body.Close()

	var data map[string]struct {
		USD float64 `json:"usd"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, false
	}

	token, ok := data[address]
	if !ok {
		return 0, false
	}

	return token.USD, true
}
