package config

type Config struct {
	AlchemyKey     string
	AlchemyURLWS   string
	AlchemyURLHTTP string
	TransferTopic  string

	CoinGeckoURL             string
	CoinGeckoAPIKey          string
	CoinGeckoPages           int
	CoinGeckoPerPage         int
	CoinGeckoTTL             int
}
