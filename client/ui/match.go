package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	shellgame "github.com/taise-hub/shellgame-cli/client"
	"github.com/taise-hub/shellgame-cli/common"
	"log"
	"sync"
	"time"
)

var muRead sync.Mutex
var muWrite sync.Mutex

type matchModel struct {
	list         list.Model
	parent       *model
	waiting      bool
	conn         *websocket.Conn
	matchingChan chan *MatchingMsg
}

func NewMatchModel() (matchModel, error) {
	l := list.New(nil, profileDelegate{}, width, 14)
	l.Title = "対戦相手を選択してください"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	mc := make(chan *MatchingMsg)

	return matchModel{list: l, waiting: false, matchingChan: mc}, nil
}

func (mm matchModel) Init() tea.Cmd {
	return nil
}

func (mm matchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
		ps, err := shellgame.GetMatchingProfiles()
		if err != nil {
			return matchModel{}, tea.Quit
		}

		var profiles []list.Item
		for _, v := range ps {
			profiles = append(profiles, Profile(*v))
		}
		mm.list.SetItems(profiles)

		// FIXME: 後で消す
		if err = shellgame.PostProfile("hoge"); err != nil {
			return matchModel{}, tea.Quit
		}

		conn, err := shellgame.ConnectMatchingRoom()
		if err != nil {
			log.Fatalf("%v\n", err.Error())
		}
		mm.conn = conn
		mm.conn.SetPongHandler(func(string) error { mm.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
		go mm.matching()
		return mm, nil
	case MatchingMsg:
		log.Printf("MatchingMsg: %+v\n", msg)
	// 受け取ったメッセージによって処理を分ける
	// 2. 対戦要求の受け取り
	// 3. 対戦要求に対する返答(DENY or ACCEPT)

	// case timeoutMsg: // 対戦要求に一定時間返答がない場合に受け取るメッセージ
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return mm, tea.Quit

		case "r":
			//対戦待ちのユーザを取得し更新する。

		case "t": // temporary ゲーム開始画面に接続する

		case "enter":
			dest, _ := mm.list.SelectedItem().(Profile)
			msg := &MatchingMsg{
				Source: dest,
				Dest: dest,
				Data: common.OFFER,
			}
			mm.matchingChan <- msg
			// 送信時に3分後にtimeoutMsgを通知する処理をgoroutineで動かす。
			// 送信後、matchModelの状態をwaitとかにしてローディング画面でも表示しとく？
			return mm, nil
		case "q":
			mm.conn.Close()
			return mm.parent, screenChange()
		}
	}
	var cmd tea.Cmd
	mm.list, cmd = mm.list.Update(msg)
	return mm, cmd
}

func (mm matchModel) View() string {
	return "\n" + mm.list.View()
}

func (mm matchModel) matching() {
	go mm.readPump()
	mm.writePump()
}

// mm.Update()から受け取ったメッセージをwebsocketに流す。
func (mm matchModel) writePump() {
	ticker := time.NewTicker(10 * time.Second)
	defer mm.conn.Close()
	for {
		select {
		case m, ok := <-mm.matchingChan:
			if !ok {
				return
			}
			if err := mm.WriteConn(m); err != nil {
				return
			}
		case <-ticker.C:
			if err := mm.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// websocketから受け取ったメッセージをmm.Update()に流す。
func (mm matchModel) readPump() {
	defer mm.conn.Close()
	p := GetProgram()
	mm.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	for {
		msg := &MatchingMsg{}
		if err := mm.ReadConn(msg); err != nil {
			return
		}
		p.Send(*msg)
	}
}

func (mm matchModel) WriteConn(msg any) error {
	defer muWrite.Unlock()
	muWrite.Lock()
	mm.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return mm.conn.WriteJSON(msg)
}

func (mm matchModel) ReadConn(msg *MatchingMsg) error {
	defer muRead.Unlock()
	muRead.Lock()
	return mm.conn.ReadJSON(msg)
}
