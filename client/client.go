package shellgame

import (
	"log"
	"bytes"
	"github.com/google/uuid"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"time"
	"encoding/json"
	"github.com/taise-hub/shellgame-cli/server/domain/model"
)

const (
	HOST = "localhost"
)

var (
	baseEndpoint  = &url.URL{Scheme: "http", Host: HOST, Path: "/"}
	profileEndpoint = &url.URL{Scheme: "http", Host: HOST, Path: "/profile"}
	shellEndpoint = &url.URL{Scheme: "ws", Host: HOST, Path: "/shell"}
)

func newClient() (*http.Client, error) {
	jar, err := getJar()
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Jar:     jar,
		Timeout: 20 * time.Second,
	}, nil
}

// シェルゲーサーバで稼働するコンテナにWebSocketを利用して接続します。
func ConnectShell() (*websocket.Conn, error) {
	var header http.Header
	jar, err := getJar()
	if err != nil {
		return nil, err
	}
	for _, cookie := range jar.Cookies(baseEndpoint) {
		header.Add("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	wsconn, _, err := websocket.DefaultDialer.Dial(shellEndpoint.String(), header)
	if err != nil {
		return nil, err
	}

	return wsconn, nil
}

// プロフィールをシェルゲーサーバに送信します。
func PostProfile(name string) error {
	id := uuid.New()
	profile := &model.Profile{ID: id.String(), Name: name}
	p, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", profileEndpoint.String(), bytes.NewBuffer(p))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client, err := newClient()
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	return nil
}