package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/taise-hub/shellgame-cli/server/domain/model"
	"github.com/taise-hub/shellgame-cli/server/usecase"
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
	b := make([]byte, req.ContentLength)
	req.Body.Read(b)
	defer req.Body.Close()
	json, err := ParseJSON(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := json["id"].(string)
	name := json["name"].(string)
	if name == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(400)
		fmt.Fprintln(w, "400 bad reuqest")
		return
	}
	sess.Values["id"] = id
	sess.Values["name"] = name
	if err := store.Save(req, w, sess); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (con *GameController) Match(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		con.getMatchingProfiles(w, req)
	default:
		http.NotFound(w, req)
	}
}

// MatchingPlayerのProfileを返すようにしたい。
func (con *GameController) getMatchingProfiles(w http.ResponseWriter, req *http.Request) {
	players := con.usecase.GetMatchingProfiles()
	RespondJSON(w, players, 200)
}

func (con *GameController) WaitMatch(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil) //NOTE: このコネクションはdomain層で利用しているためはあえて閉じてない。(domain層で閉じてる)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	sess, _ := store.Get(req, SESS_NAME)
	wc := NewWebsocketConn(conn)
	player := model.NewMatchingPlayer(sess.Values["id"].(string), sess.Values["name"].(string), wc)
	con.usecase.WaitMatch(player)
}

func ParseJSON(data []byte) (map[string]any, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func RespondJSON(w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}
