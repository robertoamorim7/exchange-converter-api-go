package main

type ParamsOld struct {
	price        float64
	toCurrencies []string
}

type Params struct {
	ToCurrencies []string `json:"toCurrencies"`
	Price        float64  `json:"price"`
}

type ConverterOutput struct {
	Message string               `json:"message"`
	Data    []map[string]float64 `json:"data"`
}

type ExchangeRate struct {
	FromCurrencyCode string `json:"1. From_Currency Code"`
	FromCurrencyName string `json:"2. From_Currency Name"`
	ToCurrencyCode   string `json:"3. To_Currency Code"`
	ToCurrencyName   string `json:"4. To_Currency Name"`
	ExchangeRate     string `json:"5. Exchange Rate"`
	LastRefreshed    string `json:"6. Last Refreshed"`
	TimeZone         string `json:"7. Time Zone"`
	BidPrice         string `json:"8. Bid Price"`
	AskPrice         string `json:"9. Ask Price"`
}

type RealtimeCurrencyExchangeRate struct {
	ExchangeRate ExchangeRate `json:"Realtime Currency Exchange Rate"`
}
