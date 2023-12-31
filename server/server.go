package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := createDB()
	if err != nil {
		log.Printf("erro ao criar banco de dados %v", err)
		panic(err)
	}

	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {

		ctxApi, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancel()
		rate, err := GetRate(ctxApi)
		if err != nil {
			log.Printf("erro ao salvar dado no banco %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Erro ao realizar consulta"))
			return
		}

		ctxBD, cancel2 := context.WithTimeout(r.Context(), 10*time.Millisecond)
		defer cancel2()
		err = NewRepository(db).Create(ctxBD, rate)
		if err != nil {
			log.Printf("Erro ao salvar dado no banco %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Erro ao realizar consulta"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Bid string `json:"bid"`
		}{Bid: rate.DollarExchangeRate.Bid})

	})

	http.ListenAndServe(":8080", mux)

}

/**
Struct
**/

type Rate struct {
	DollarExchangeRate struct {
		Code      string `json:"code"`
		Codein    string `json:"codein"`
		Name      string `json:"name"`
		High      string `json:"high"`
		Low       string `json:"low"`
		VarBid    string `json:"varBid"`
		PctChange string `json:"pctChange"`
		Bid       string `json:"bid"`
		Ask       string `json:"ask"`
	} `json:"USDBRL"`
}

/**
Codigo do Repository
**/

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (repository *Repository) Create(ctx context.Context, entity *Rate) error {
	select {
	case <-ctx.Done():
		log.Println("Timeout na persistência no banco de dados")
		return fmt.Errorf("timeout ao salvar dados")
	default:
		query, err := repository.db.Prepare("INSERT INTO dollar_exchange_rate" +
			"(code, codein, name, high, low, varBid, pctChange, bid, ask) " +
			"VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer query.Close()

		_, err = query.Exec(entity.DollarExchangeRate.Code, entity.DollarExchangeRate.Codein, entity.DollarExchangeRate.Name, entity.DollarExchangeRate.High, entity.DollarExchangeRate.Low, entity.DollarExchangeRate.VarBid, entity.DollarExchangeRate.PctChange, entity.DollarExchangeRate.Bid, entity.DollarExchangeRate.Ask)
		if err != nil {
			log.Printf("Erro ao salvar dado no banco %v", err)
			return err
		}

		return nil
	}
}

/**
API
**/

func GetRate(ctx context.Context) (*Rate, error) {

	select {
	case <-ctx.Done():
		log.Println("Timeout ao consultar api")
		return nil, fmt.Errorf("timeout ao consultar api")
	default:
		res, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")

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

		var data Rate
		err = json.Unmarshal(result, &data)
		if err != nil {
			log.Printf("Erro ao realizar parse da req %v", err)
			return nil, err
		}
		return &data, nil
	}

}

func createDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./data.db")

	if err != nil {
		return nil, err
	}

	createTableSQL := `
		CREATE TABLE dollar_exchange_rate (
			id SERIAL PRIMARY KEY,
			code VARCHAR(10) NOT NULL,
			codein VARCHAR(10) NOT NULL,
			name VARCHAR(255) NOT NULL,
			high FLOAT NOT NULL,
			low FLOAT NOT NULL,
			varBid FLOAT NOT NULL,
			pctChange FLOAT NOT NULL,
			bid FLOAT NOT NULL,
			ask FLOAT NOT NULL
		);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}
	return db, nil

}
