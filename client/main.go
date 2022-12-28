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

type CotacaoDados struct {
	Cotacao string `json:"bid"`
}

func main() {
	// Set a timeout of 300 ms
	// In case no response is returned, it returns an error of context deadline exceeded
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Since we are using context, we need to make the request with context
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	// We execute the request we created
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	// Then we get the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var response CotacaoDados

	// Parse the JSON to struct
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// Create a file that is going to store the exchange rate... in this case, USDolar -> BRReal
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	// Write the response to the file
	_, err = f.Write([]byte("Dólar:" + response.Cotacao))
	if err != nil {
		panic(err)
	}
	f.Close()

	fmt.Println(response.Cotacao)
}
