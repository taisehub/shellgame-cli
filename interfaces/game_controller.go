package interfaces

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/taise-hub/shellgame-cli/usecase"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type GameController struct {
	usecase usecase.GameUsecase
}

func NewGameController(usecase usecase.GameUsecase) *GameController {
	return &GameController{
		usecase: usecase,
	}
}

func (con *GameController) Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello, %v\n", req.FormValue("name"))
}

func (con *GameController) HandleWebsocket(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	defer conn.Close()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if err = con.usecase.Start(conn.UnderlyingConn()); err != nil {
		fmt.Fprintf(w, "error")
	}
}
