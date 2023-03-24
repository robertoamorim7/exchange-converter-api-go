package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func asyncConvert(fromCurrency string, currency string, price float64, wg *sync.WaitGroup, results chan float64) {
	defer wg.Done()

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s&apikey=%s", fromCurrency, currency, API_KEY)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var data RealtimeCurrencyExchangeRate
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err.Error())
	}

	if ok := data.ExchangeRate.ExchangeRate; ok == "" {
		fmt.Println("Realtime Currency Exchange Rate not in response") //change
	}

	exchangeRate, err := strconv.ParseFloat(data.ExchangeRate.ExchangeRate, 64)
	if err != nil {
		fmt.Println(err.Error())
	}

	results <- exchangeRate * price
}
