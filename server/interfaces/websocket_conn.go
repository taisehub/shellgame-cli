package interfaces

import (
	"github.com/gorilla/websocket"
	"github.com/taise-hub/shellgame-cli/server/domain/model"
	"sync"
	"time"
)

const (
	writeWait      = 10 * time.Second
	readWait       = 60 * time.Second
	maxMessageSize = 512
)

type WebsocketConn struct {
	*websocket.Conn
}

func NewWebsocketConn(conn *websocket.Conn) *WebsocketConn {
	conn.SetReadLimit(maxMessageSize)
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	conn.SetReadDeadline(time.Now().Add(readWait))
	conn.SetPingHandler(func(string) error { 
		conn.SetReadDeadline(time.Now().Add(readWait));
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err :=conn.WriteMessage(websocket.PongMessage, nil); err != nil {
			return err
		}
		return nil
	 })
	return &WebsocketConn{conn}
}

// コネクションをcloseする前にCloseを通知するメッセージを送信することにする。
// shellgame-clientの実装によっては不要になる可能性あり。
func (wc *WebsocketConn) Close() error {
	mu := sync.Mutex{}
	defer mu.Unlock()
	mu.Lock()
	wc.WriteMessage(websocket.CloseMessage, []byte{})
	return wc.Conn.Close()
}

func (wc *WebsocketConn) Read(msg model.Message) error {
	mu := sync.Mutex{}
	defer mu.Unlock()
	mu.Lock()
	return wc.ReadJSON(msg)
}

func (wc *WebsocketConn) Write(msg model.Message) error {
	mu := sync.Mutex{}
	defer mu.Unlock()
	mu.Lock()
	return wc.WriteJSON(msg)
}
