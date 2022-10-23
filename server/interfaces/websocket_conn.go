package interfaces

import (
	"github.com/gorilla/websocket"
	"github.com/taise-hub/shellgame-cli/common"
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
	muRead           sync.Mutex
	muWrite          sync.Mutex
}

func NewWebsocketConn(conn *websocket.Conn) *WebsocketConn {
	conn.SetReadLimit(maxMessageSize)
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	conn.SetReadDeadline(time.Now().Add(readWait))
	conn.SetPingHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(readWait))
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
			return err
		}
		return nil
	})
	return &WebsocketConn{conn, sync.Mutex{}, sync.Mutex{}}
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

func (wc *WebsocketConn) Read(msg common.Message) error {
	defer wc.muRead.Unlock()
	wc.muRead.Lock()
	return wc.ReadJSON(msg)
}

func (wc *WebsocketConn) Write(msg common.Message) error {
	defer wc.muWrite.Unlock()
	wc.muWrite.Lock()
	return wc.WriteJSON(msg)
}
