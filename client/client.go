package shellgame

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/taise-hub/shellgame-cli/common"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	HOST = "localhost"
)

var (
	baseEndpoint     = &url.URL{Scheme: "http", Host: HOST, Path: "/"}
	profileEndpoint  = &url.URL{Scheme: "http", Host: HOST, Path: "/profiles"}
	playersEndpoint  = &url.URL{Scheme: "http", Host: HOST, Path: "/players"}
	shellEndpoint    = &url.URL{Scheme: "ws", Host: HOST, Path: "/shell"}
	matchingEndpoint = &url.URL{Scheme: "ws", Host: HOST, Path: "/waitmatch"}
)

// シェルゲーサーバで稼働するコンテナにWebSocketを利用して接続する。
func ConnectShell() (*websocket.Conn, error) {
	jar, err := getJar()
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	for _, cookie := range jar.Cookies(baseEndpoint) {
		header.Add("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}

	wsconn, _, err := websocket.DefaultDialer.Dial(shellEndpoint.String(), header)
	if err != nil {
		return nil, err
	}
	return wsconn, nil
}

// シェルゲーサーバで稼働するマッチングルームにWebSocketを利用して接続する。
func ConnectMatchingRoom() (*websocket.Conn, error) {
	jar, err := getJar()
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	for _, cookie := range jar.Cookies(baseEndpoint) {
		header.Add("Cookie", fmt.Sprintf("%s", cookie))
	}

	wsconn, _, err := websocket.DefaultDialer.Dial(matchingEndpoint.String(), header)
	if err != nil {
		return nil, err
	}
	return wsconn, nil
}

// シェルゲーサーバにプレイヤー名を登録する。
func PostProfile(name string) error {
	id := uuid.New()
	profile := &common.Profile{ID: id.String(), Name: name}
	p, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	resp, err := http.Post(profileEndpoint.String(), "application/json", bytes.NewBuffer(p))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", body)
	}
	jar, err := getJar()
	if err != nil {
		return err
	}
	jar.SetCookies(baseEndpoint, resp.Cookies())
	return nil
}

// シェルゲーサーバから対戦待ちユーザを取得する
func GetMatchingProfiles() ([]*common.Profile, error) {
	resp, err := http.Get(playersEndpoint.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var profiles []*common.Profile
	if err := json.Unmarshal(body, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}
