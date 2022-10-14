package ui

import (
	"log"
	"time"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"sync"
	shellgame "github.com/taise-hub/shellgame-cli/client"
)

var mu sync.Mutex

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

	err := shellgame.PostProfile("hoge")
	if err != nil {
		return matchModel{}, err
	}
	players, err := shellgame.GetMatchingPlayers()
	if err != nil {
		return matchModel{}, err
	}

	var profiles []list.Item
	for _, v := range players {
		profiles = append(profiles, Profile{v.Profile.ID, v.Profile.Name})
	}
	// FIXME: 下2行いつか消す
	profiles = append(profiles, Profile{"1", "bob"})
	profiles = append(profiles, Profile{"2", "alice"})
	l.SetItems(profiles)

	return matchModel{list: l, waiting: false, matchingChan: mc}, nil
}

func (mm matchModel) Init() tea.Cmd {
	// FIXME: screenChangeMsgを受け取ったた時にコネクション張る方がいいかも？
	conn, err := shellgame.ConnectMatchingRoom()
	if err != nil {
		log.Fatalf("%v\n", err.Error())
	}
	mm.conn = conn
	go mm.matching()
	return nil
}

func (mm matchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MatchingMsg:
	// 受け取ったメッセージによって処理を分ける
	// 2. 対戦要求の受け取り
	// 3. 対戦要求に対する返答(DENY or ACCEPT)
	// case timeoutMsg: // 対戦要求に一定時間返答がない場合に受け取るメッセージ
	case screenChangeMsg:
		return mm, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return mm, tea.Quit

		case "r":
			//対戦待ちのユーザを取得し更新する。

		case "enter":
			// チャネルに書き込んで、writePump()にデータを流せばOK
			msg := &MatchingMsg{}
			dest, _ := mm.list.SelectedItem().(Profile)
			msg.Dest = dest
			mm.matchingChan <- msg
			// 送信時に3分後にtimeoutMsgを通知する処理をgoroutineで動かす。
			// 送信後、matchModelの状態をwaitとかにしてローディング画面でも表示しとく？
			return mm, nil
		case "q":
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
	defer mm.conn.Close()
	for {
		m, ok := <- mm.matchingChan
		if !ok {
			return
		}
		if err := mm.WriteConn(m); err != nil {
			return
		}
	}
}

// websocketから受け取ったメッセージをmm.Update()に流す。
func (mm matchModel) readPump() {
	defer mm.conn.Close()
	p := GetProgram()
	for { 
		var msg *MatchingMsg
		if err := mm.ReadConn(msg); err != nil {
			return
		}
		p.Send(msg)
	}
}

func (mm matchModel) WriteConn(msg *MatchingMsg) error {
	defer mu.Unlock()
	mu.Lock()
	mm.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return mm.conn.WriteJSON(msg)
}

func (mm matchModel) ReadConn(msg *MatchingMsg) error {
	defer mu.Unlock()
	mu.Lock()
	mm.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return mm.conn.ReadJSON(msg)
}