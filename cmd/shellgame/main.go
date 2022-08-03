package main

import (
	"github.com/taise-hub/shellgame-cli/infrastructure"
	"github.com/taise-hub/shellgame-cli/interfaces"
	"github.com/taise-hub/shellgame-cli/usecase"
	"log"
	"net/http"
)

func main() {
	containerHandler, err := infrastructure.NewContainerHandler()
	if err != nil {
		log.Fatal(err)
		return
	}
	consoleRepo := interfaces.NewContainerRepository(containerHandler)
	gameUsecase := usecase.NewGameInteractor(consoleRepo)
	gameController := interfaces.NewGameController(gameUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", gameController.Hello)
	mux.HandleFunc("/ws", gameController.HandleWebsocket)

	log.Println("[+] Start listening.")
	http.ListenAndServe(":80", mux)
}
