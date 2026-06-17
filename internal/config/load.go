package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func getEnv(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AlchemyKey:     getEnv("AlchemyKey", ""),
		AlchemyURLWS:   getEnv("AlchemyURLWS", ""),
		AlchemyURLHTTP: getEnv("AlchemyURLHTTP", ""),
		TransferTopic:  getEnv("TransferTopic", ""),

		CoinGeckoURL:     getEnv("COINGECKO_URL", "https://api.coingecko.com/"),
		CoinGeckoAPIKey:  getEnv("COINGECKO_API_KEY", ""),
		CoinGeckoPages:   getEnvInt("COINGECKO_PAGES", 4),
		CoinGeckoPerPage: getEnvInt("COINGECKO_PER_PAGE", 200),
		CoinGeckoTTL:     getEnvInt("COINGECKO_TTL", 60),
	}
}
