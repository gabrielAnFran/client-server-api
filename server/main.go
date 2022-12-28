package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Cotacao struct {
	Cotacao CotacaoDados `json:"USDBRL"`
}

type CotacaoDados struct {
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

	// Open conection with db
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/clientserver")
	if err != nil {
		return nil, err
	}

	// Call function that is going to store the current exchange rate
	// Using a go routime because the client wants to get the exchange rate
	// and for them it doesn't matter if the server was able to persist the data
	// since it is an internal controll
	// it would be interesting using a monitoring system to get notify if an error occurs in that function
	// for example Sentry
	go salvarCotacaoAtual(db, response.Cotacao.Bid, response.Cotacao.Timestamp)

	return &response, err

}

func salvarCotacaoAtual(db *sql.DB, bid, timestamp string) {
	// Set context with 10ms of timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Prepare statemnte
	// Good practice for avoiding sql injection
	stmt, err := db.Prepare("insert into cotacoes(id, timestamp, bid) values(?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		// Send error to monitoring system
	}
	defer stmt.Close()

	// Exec insert using the pre-defined context
	_, err = stmt.ExecContext(ctx, uuid.New().String(), timestamp, bid)
	if err != nil {
		fmt.Println(err)
		// Send error to monitoring system
	}

}
