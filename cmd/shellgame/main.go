package main

import (
	"github.com/taise-hub/shellgame-cli/interfaces"
	"log"
	"net/http"
)

func main() {
	gameController := interfaces.NewGameController()
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", gameController.Hello)
	mux.HandleFunc("/ws", gameController.HandleWebsocket)

	log.Println("[+] Start listening.")
	http.ListenAndServe(":", mux)
}
