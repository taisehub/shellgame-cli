package interfaces

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
	"github.com/taise-hub/shellgame-cli/usecase"
	"github.com/taise-hub/shellgame-cli/domain/model"
	"net/http"
)

const (
	SESS_NAME = "shellgame-sess"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
)

type GameController struct {
	usecase usecase.GameUsecase
}

func NewGameController(usecase usecase.GameUsecase) *GameController {
	return &GameController{
		usecase: usecase,
	}
}

// 対戦開始時にクライアントから呼び出される予定
// websocketを用いてクライアントをシェルに接続する
func (con *GameController) Start(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	defer conn.Close()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err = con.usecase.Start(conn.UnderlyingConn()); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (con *GameController) Profile(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		con.saveProfile(w, req)
	default:
		http.NotFound(w, req)
	}
}

func (con *GameController) saveProfile(w http.ResponseWriter, req *http.Request) {
	sess, _ := store.Get(req, SESS_NAME)
	name := req.FormValue("name")
	if name == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(400)
		fmt.Fprintln(w, "400 bad reuqest")
		return
	}
	// もうちょっとスマートにidを生成をしたい
	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)
  	sess.Values["id"] = random.Uint32()
	// ------------
  	sess.Values["name"] = name
  	if err := store.Save(req, w, sess); err != nil {
  	  http.Error(w, err.Error(), http.StatusInternalServerError)
  	  return
  	}
}

func (con *GameController) WaitMatch(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil) //NOTE: このコネクションはdomain層で利用しているためはあえて閉じてない。(domain層で閉じてる)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sess, _ := store.Get(req, SESS_NAME)
	player := model.NewMatchingPlayer(sess.Values["id"].(uint32), sess.Values["name"].(string), &WebsocketConn{conn, sync.Mutex{}})
	con.usecase.WaitMatch(player)
}
