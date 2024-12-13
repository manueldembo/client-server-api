package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/cotacao", ExchangeRateHandler)
	http.ListenAndServe(":8080", nil)
}

func ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Accept", "application/json")

	e, err := getExchangeRate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(e)
}

func getExchangeRate() (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var exchangeRate ExchangeRate
	err = json.NewDecoder(res.Body).Decode(&exchangeRate)
	if err != nil {
		return nil, err
	}

	return &exchangeRate.USDBRL.Bid, nil
}

type ExchangeRate struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}
