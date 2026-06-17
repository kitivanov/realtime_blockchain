package model

type TokenMeta struct {
	Address  string
	Symbol   string
	Name     string
	Decimals int
}

type TokenList struct {
	Name    string        `json:"name"`
	LogoURI string        `json:"logoURI"`
	Tokens  []TokenRecord `json:"tokens"`
}

type TokenRecord struct {
	ChainId  int    `json:"chainId"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	LogoURI  string `json:"logoURI"`
}
