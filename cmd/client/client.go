package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		fmt.Printf("Error: %s\n", string(body))
		return
	}

	var exchangeRate string
	err = json.NewDecoder(res.Body).Decode(&exchangeRate)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	saveExchangeRate(exchangeRate)
}

func saveExchangeRate(value string) error {
	f, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("Dolar: %s\n", value))
	if err != nil {
		return err
	}

	println("Content added successfully!")
	return nil
}
