package exporter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/D3vR4pt0rs/logger"
)

type client struct {
	exporterAddress string
	httpClient      *http.Client
}

func New() *client {
	address := os.Getenv("EXPORTER_ADDRESS")
	return &client{
		exporterAddress: address,
		httpClient:      &http.Client{},
	}
}

type Ticker struct {
	Symbol    string    `json:"symbol"`
	LastPrice float64   `json:"last_price"`
	Timestamp time.Time `json:"timestamp"`
}

type Response struct {
	Data Ticker `json:"data"`
}

func (c client) GetTickerPrice(ticker string) (float64, error) {
	url := fmt.Sprintf("http://%s/api/ticker/%s", c.exporterAddress, ticker)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		logger.Error.Println(err.Error())
		return 0, err
	}

	logger.Info.Println(url, resp.StatusCode)
	defer resp.Body.Close()

	var result Response

	json.NewDecoder(resp.Body).Decode(&result)
	logger.Info.Println("Ticker information: ", result)

	return result.Data.LastPrice, nil
}
