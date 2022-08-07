package main

import (
	"github.com/taise-hub/shellgame-cli/domain/model"
	"github.com/taise-hub/shellgame-cli/domain/service"
	"github.com/taise-hub/shellgame-cli/infrastructure"
	"github.com/taise-hub/shellgame-cli/interfaces"
	"github.com/taise-hub/shellgame-cli/interfaces/redis"
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
	redisHandler := infrastructure.NewRedisHandler()
	matchingRoomRepo := redis.NewMatchingRoomRepository(redisHandler)
	matchService := service.NewMatchService(matchingRoomRepo)
	gameUsecase := usecase.NewGameInteractor(consoleRepo, matchService)
	gameController := interfaces.NewGameController(gameUsecase)

	go model.GetMatchingRoom().Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/profile", gameController.Profile)
	mux.HandleFunc("/match/wait", gameController.WaitMatch)
	mux.HandleFunc("/ws", gameController.Start)

	log.Println("[+] Start listening.")
	http.ListenAndServe(":80", mux)
}
