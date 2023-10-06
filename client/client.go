package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	ctxApi, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	rate, err := GetBid(ctxApi)
	if err != nil {
		log.Printf("Erro ao realizar consulta %v", err)
		panic(err)
	}

	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte(fmt.Sprintf("Dólar:%s", rate.Bid)))
	if err != nil {
		panic(err)
	}
	f.Close()

}

type Bid struct {
	Bid string `json:"bid"`
}

/**
API
**/

func GetBid(ctx context.Context) (*Bid, error) {

	select {
	case <-ctx.Done():
		log.Println("Timeout ao consultar api")
		return nil, fmt.Errorf("timeout ao consultar api")
	default:
		res, err := http.Get("http://localhost:8080/cotacao")

		if err != nil {
			log.Printf("Erro ao realizar consulta %v", err)
			return nil, err
		}

		defer res.Body.Close()
		result, err := io.ReadAll(res.Body)

		if err != nil {
			log.Printf("Erro ao realizar parse da req %v", err)
			return nil, err

		}

		var data Bid
		err = json.Unmarshal(result, &data)
		if err != nil {
			log.Printf("Erro ao realizar parse da req %v", err)
			return nil, err
		}
		return &data, nil
	}

}
