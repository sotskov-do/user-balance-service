package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func currencyConversion(curr string, v float64) (float64, error) {
	url := "http://api.exchangeratesapi.io/v1/latest?access_key=15934be9b0a839cc7998faeb9e3babf9"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err) // TODO обработать
	}

	currencies := map[string]interface{}{}
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		panic(err) // TODO обработать
	}
	c, ok := currencies["rates"].(map[string]interface{})[curr].(float64)
	if !ok {
		return 0, fmt.Errorf("can't find such currency. please use another one")
	}

	result := v / currencies["rates"].(map[string]interface{})["RUB"].(float64) * c

	return result, nil
}
