package main

import (
	"log"
	"net/http"
	"time"
)

func start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		select {
		case <-time.After(5 * time.Second):
			log.Println("Request finalizada")
			w.Write([]byte("Request finalizada"))

		case <-ctx.Done():
			log.Println("Request timeout")
			w.Write([]byte("Request Cancelada"))

		}
	})

	http.ListenAndServe(":8080", mux)

}
