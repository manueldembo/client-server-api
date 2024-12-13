package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Accept", "application/json")

		e, err := getExchangeRate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = saveExchangeRate(db, e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(e)
	})
	http.ListenAndServe(":8080", nil)
}

func getExchangeRate() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var exchangeRate ExchangeRateDTO
	err = json.NewDecoder(res.Body).Decode(&exchangeRate)
	if err != nil {
		return "", err
	}

	return exchangeRate.USDBRL.Bid, nil
}

func saveExchangeRate(db *sql.DB, value string) error {
	if value == "" {
		return errors.New("invalid exchange rate")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	id := uuid.New().String()
	_, err := db.ExecContext(ctx, "INSERT INTO exchange_rates (id, value) VALUES ($1, $2)", id, value)
	if err != nil {
		return err
	}

	return nil
}

type ExchangeRateDTO struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}
