package interfaces

import (
	"github.com/gorilla/websocket"
	"github.com/taise-hub/shellgame-cli/domain/model"
	"sync"
	"time"
)

const (
	writeWait      = 10 * time.Second
	readWait       = 10 * time.Second
	maxMessageSize = 512
)

type WebsocketConn struct {
	*websocket.Conn
	sync.Mutex
}

// コネクションをcloseする前にCloseを通知するメッセージを送信することにする。
// shellgame-clientの実装によっては不要になる可能性あり。
func (wc *WebsocketConn) Close() error {
	defer wc.Unlock()
	wc.Lock()
	wc.WriteMessage(websocket.CloseMessage, []byte{})
	return wc.Conn.Close()
}

func (wc *WebsocketConn) Read(msg model.Message) error {
	defer wc.Unlock()
	wc.Lock()
	wc.SetReadLimit(maxMessageSize)
	wc.SetReadDeadline(time.Now().Add(readWait))
	return wc.ReadJSON(msg)
}

func (wc *WebsocketConn) Write(msg model.Message) error {
	defer wc.Unlock()
	wc.Lock()
	wc.SetWriteDeadline(time.Now().Add(writeWait))
	return wc.WriteJSON(msg)
}
