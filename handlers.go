package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var API_KEY string = os.Getenv("API_KEY")

func syncConverter(c *fiber.Ctx) error {
	// controller layer
	fromCurrency := c.Params("from_currency")
	toCurrenciesString := c.Query("to_currencies")
	priceString := c.Query("price")
	price, err := strconv.ParseFloat(priceString, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	toCurrencies := strings.Split(toCurrenciesString, ",")

	// service layer
	var url string
	var exchangeRates []float64

	for _, currency := range toCurrencies {
		url = fmt.Sprintf("https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s&apikey=%s", fromCurrency, currency, API_KEY)
		resp, err := http.Get(url)
		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		var data RealtimeCurrencyExchangeRate
		err = json.Unmarshal(body, &data)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if ok := data.ExchangeRate.ExchangeRate; ok == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "Realtime Currency Exchange Rate not in response")
		}

		exchangeRate, err := strconv.ParseFloat(data.ExchangeRate.ExchangeRate, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}

		exchangeRates = append(exchangeRates, exchangeRate*price)
	}

	return c.JSON(exchangeRates)
}

func asyncConverter(c *fiber.Ctx) error {
	// controller layer
	fromCurrency := c.Params("from_currency")
	toCurrencies := c.Query("to_currencies")
	toCurrenciesList := strings.Split(toCurrencies, ",")
	price, err := strconv.ParseFloat(c.Query("price"), 64)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	err = validateParams(fromCurrency, toCurrencies, price)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	// service layer
	var wg sync.WaitGroup
	results := make(chan float64, len(toCurrenciesList))

	for _, currency := range toCurrenciesList {
		wg.Add(1)
		go asyncConvert(fromCurrency, currency, price, &wg, results)
	}

	wg.Wait()

	var convertedPrices []float64

	for i := 0; i < len(toCurrenciesList); i++ {
		convertedPrices = append(convertedPrices, <-results)
	}

	return c.JSON(convertedPrices)
}

func asyncConverterV2(c *fiber.Ctx) error {
	// controller layer
	var params Params
	fromCurrency := c.Params("from_currency")
	bodyBytes := c.Body()

	err := json.Unmarshal(bodyBytes, &params)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	err = validateParamsV2(params)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	// service layer
	var response []map[string]float64
	var wg sync.WaitGroup
	results := make(chan float64, len(params.ToCurrencies))

	for _, currency := range params.ToCurrencies {
		wg.Add(1)
		go asyncConvert(fromCurrency, currency, params.Price, &wg, results)
	}

	wg.Wait()

	for i := 0; i < len(params.ToCurrencies); i++ {
		result := <-results
		response = append(response, map[string]float64{
			params.ToCurrencies[i]: result,
		})
	}

	output := ConverterOutput{
		Message: "it worked!",
		Data:    response,
	}

	return c.JSON(output)
}
