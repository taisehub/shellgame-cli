package shellgame

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"time"
)

const (
	HOST = "localhost"
)

var (
	BaseEndpoint  = &url.URL{Scheme: "http", Host: HOST, Path: "/"}
	ShellEndpoint = &url.URL{Scheme: "ws", Host: HOST, Path: "/shell"}
)

func NewClient() (*http.Client, error) {
	jar, err := getJar()
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Jar:     jar,
		Timeout: 20 * time.Second,
	}, nil
}

// Connect() connect to container hosted by shellgame server using websocket.
func ConnectShell() (*websocket.Conn, error) {
	var header http.Header
	jar, err := getJar()
	if err != nil {
		return nil, err
	}
	for _, cookie := range jar.Cookies(BaseEndpoint) {
		header.Add("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	wsconn, _, err := websocket.DefaultDialer.Dial(ShellEndpoint.String(), header)
	if err != nil {
		return nil, err
	}

	return wsconn, nil
}
