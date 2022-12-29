package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	Cotacao CotacaoDados `json:"USDBRL"`
}

type CotacaoDados struct {
	Id         int    `gorm:"primaryKey"`
	Ask        string `json:"ask"`
	Bid        string `json:"bid"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	CreateDate string `json:"create_date"`
	High       string `json:"high"`
	Low        string `json:"low"`
	Name       string `json:"name"`
	PctChange  string `json:"pctChange"`
	Timestamp  string `json:"timestamp"`
	VarBid     string `json:"varBid"`
}

type CotacaoDadosPersistir struct {
	Id        int    `gorm:"primaryKey"`
	Bid       string `json:"bid"`
	Timestamp string `json:"timestamp"`
	DataInc   time.Time
}

func main() {
	fmt.Println("Server running...")
	mux := http.NewServeMux()

	mux.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", mux)

}

func handler(w http.ResponseWriter, r *http.Request) {

	cotacao, err := cotacaoDolarBuscar()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(cotacao.Cotacao)

}

func cotacaoDolarBuscar() (*Cotacao, error) {
	// Set a timeout of 200ms
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Prepare the request with context
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	// Get data from the res.Body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response Cotacao

	// Parse it to struct
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Call function that is going to store the current exchange rate
	// Using a go routime because the client wants to get the exchange rate
	// and for them it doesn't matter if the server was able to persist the data
	// since it is an internal controll
	// it would be interesting using a monitoring system to get notify if an error occurs in that function
	// for example Sentry
	go salvarCotacaoAtual(response.Cotacao)

	return &response, err

}

func salvarCotacaoAtual(cotacao CotacaoDados) {
	// Set context with 10ms of timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	err = db.AutoMigrate(&CotacaoDadosPersistir{})
	if err != nil {
		fmt.Println(err)
	}

	cotacaoInserir := &CotacaoDadosPersistir{
		Bid:       cotacao.Bid,
		Timestamp: cotacao.Timestamp,
		DataInc:   time.Now(),
	}
	err = db.WithContext(ctx).
		Create(&cotacaoInserir).Error
	if err != nil {
		fmt.Println(err)
	}

}
