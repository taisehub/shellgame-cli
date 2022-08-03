package interfaces

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/taise-hub/shellgame-cli/usecase"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type GameController struct {
	usecase *usecase.GameUsecase
}

func NewGameController() *GameController {
	return &GameController{
		usecase: usecase.NewGameUsecase(),
	}
}

func (con *GameController) Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello, %v\n", req.FormValue("name"))
}

func (con *GameController) HandleWebsocket(w http.ResponseWriter, req *http.Request) {
	_, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
