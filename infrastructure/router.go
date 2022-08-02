package infrastructure

import (
	"log"
	"net/http"
)

var (
	port = "80"
)

func Route() error {
	mux := http.NewServeMux()
	mux.Handle("/hello", nil)
	mux.Handle("/ws", nil)
	log.Println("[+] Start listening.")
	return http.ListenAndServe(":" + port, mux)
}