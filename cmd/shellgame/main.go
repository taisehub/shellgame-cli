package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/hello", nil)
	mux.Handle("/ws", nil)

	log.Println("[+] Start listening.")
	http.ListenAndServe(":", mux)
}
