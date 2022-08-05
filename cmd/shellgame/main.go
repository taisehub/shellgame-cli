package main

import (
	"github.com/taise-hub/shellgame-cli/infrastructure"
	"github.com/taise-hub/shellgame-cli/interfaces"
	"github.com/taise-hub/shellgame-cli/interfaces/redis"
	"github.com/taise-hub/shellgame-cli/usecase"
	"github.com/taise-hub/shellgame-cli/domain/service"
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
	redisHandler := infrastructure.NewRedisHandler()
	matchingRoomRepo := redis.NewMatchingRoomRepository(redisHandler)
	matchService := service.NewMatchService(matchingRoomRepo)
	gameUsecase := usecase.NewGameInteractor(consoleRepo, matchService)
	gameController := interfaces.NewGameController(gameUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", gameController.Start)

	log.Println("[+] Start listening.")
	http.ListenAndServe(":80", mux)
}
