package main

import (
	"errors"
	"regexp"
)

func validateParams(fromCurrency string, toCurrencies string, price float64) error {
	var err error

	err = validatePathParam(fromCurrency)
	if err != nil {
		return err
	}
	err = validateToCurrenciesParam(toCurrencies)
	if err != nil {
		return err
	}
	// validate price param
	if price < 0 {
		return errors.New("invalid query parameter")
	}

	return nil
}

func validatePathParam(param string) error {
	regex, err := regexp.Compile("^[A-Z]{3}$")
	if err != nil {
		return errors.New("invalid path parameter")
	}
	if validParam := regex.MatchString(param); !validParam {
		return errors.New("invalid path parameter")
	}

	if len(param) > 3 {
		return errors.New("invalid path parameter")
	}

	return nil
}

func validateToCurrenciesParam(param string) error {
	regex, err := regexp.Compile("^[A-Z]{3}(,[A-Z]{3})$")
	if err != nil {
		return errors.New("invalid query parameter")
	}
	if validParam := regex.MatchString(param); !validParam {
		return errors.New("invalid query parameter")
	}

	if len(param) > 50 {
		return errors.New("too much parameters")
	}

	return nil
}

func validateParamsV2(params Params) error {
	// validate to currencies
	regex, err := regexp.Compile("^[A-Z]{3}$")
	if err != nil {
		return errors.New("invalid query parameter")
	}
	for _, toCurrency := range params.ToCurrencies {
		if validParam := regex.MatchString(toCurrency); !validParam {
			return errors.New("invalid query parameter")
		}
	}

	//validate price
	if params.Price < 0 {
		return errors.New("invalid query parameter")
	}

	return nil
}
