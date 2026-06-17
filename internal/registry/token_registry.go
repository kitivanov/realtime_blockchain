package registry

import (
	"encoding/json"
	"os"
	"realtime_blockchain/internal/model"
	"strings"
	"sync"
)

type TokenRegistry struct {
	mu sync.RWMutex

	tokens map[string]*model.TokenMeta
}

func NewTokenRegistryFromFile(path string) (*TokenRegistry, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var list model.TokenList
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}

	reg := &TokenRegistry{
		tokens: make(map[string]*model.TokenMeta),
	}

	for _, t := range list.Tokens {
		addr := normalizeAddr(t.Address)

		reg.tokens[addr] = &model.TokenMeta{
			Address:  addr,
			Symbol:   t.Symbol,
			Name:     t.Name,
			Decimals: t.Decimals,
		}
	}

	return reg, nil
}

func (r *TokenRegistry) Get(address string) (*model.TokenMeta, bool) {
	addr := normalizeAddr(address)

	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tokens[addr]
	return t, ok
}

func (r *TokenRegistry) Upsert(meta *model.TokenMeta) {
	if meta == nil {
		return
	}

	addr := normalizeAddr(meta.Address)

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokens[addr] = meta
}

func normalizeAddr(addr string) string {
	return strings.ToLower(addr)
}
