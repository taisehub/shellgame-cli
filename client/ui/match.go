package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	shellgame "github.com/taise-hub/shellgame-cli/client"
	"github.com/taise-hub/shellgame-cli/common"
	"time"
)

type matchModel struct {
	list         list.Model
	parent       *topModel
	screen       screen
	choice       choiceModel
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

	choice := NewChoiceModel()

	return matchModel{list: l, screen: "", choice: choice, matchingChan: mc}, nil
}

func (mm matchModel) Init() tea.Cmd {
	return nil
}

func (mm matchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch mm.screen {
	case "choice":
		return mm.choice.Update(msg, mm)
	default:
		return mm.update(msg)
	}
}

func (mm matchModel) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case screenChangeMsg:
		return mm.screenChangeHandler(msg)
	case MatchingMsg:
		return mm.matchingMsgHandler(msg)
	// case timeoutMsg: // 対戦要求に一定時間返答がない場合に受け取るメッセージ
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return mm, tea.Quit
		case "enter":
			dest, _ := mm.list.SelectedItem().(Profile)
			if dest.ID == "" {
				return mm, nil
			}
			msg := &MatchingMsg{
				Source: dest,
				Dest:   dest,
				Data:   common.OFFER,
			}
			mm.matchingChan <- msg
			// 送信時に3分後にtimeoutMsgを通知する処理をgoroutineで動かす。
			// 送信後、matchModelの状態をwaitとかにしてローディング画面でも表示しとく？
			return mm, nil
		case "q":
			mm.conn.Close()
			return mm.parent, screenChange("match")
		}
	}
	var cmd tea.Cmd
	mm.list, cmd = mm.list.Update(msg)
	return mm, cmd
}

func (mm matchModel) View() string {
	switch mm.screen {
	case "choice":
		return mm.choice.View()
	default:
		return "\n" + mm.list.View()
	}
}

func (mm matchModel) screenChangeHandler(msg screenChangeMsg) (tea.Model, tea.Cmd) {
	switch msg {
	case "top": // TOP画面からの遷移。現在対戦待ちのPlayerを取得し、webosocketでコネクションを生成する。
		if err := mm.updateProfiles(); err != nil {
			return matchModel{}, tea.Quit
		}
		if err := mm.createConn(); err != nil {
			return matchModel{}, tea.Quit
		}
		go mm.matching()
		return mm, nil
	case "choice": // 対戦要求などの回答画面からの遷移。現在対戦待ちのPlayerを更新する。
		if err := mm.updateProfiles(); err != nil {
			return matchModel{}, tea.Quit
		}
		return mm, nil
	}
	return mm, nil
}

func (mm matchModel) matchingMsgHandler(msg MatchingMsg) (tea.Model, tea.Cmd) {
	switch msg.Data {
	case common.OFFER:
		mm.screen = "choice"
		return mm, screenChange("match")
	case common.CANCEL_OFFER:
	case common.ACCEPT:
	case common.DENY:
	case common.ERROR:
	case common.JOIN:
		mm.appendProfile(msg.Source)
		return mm, nil
	case common.LEAVE:
		mm.removeProfile(msg.Source)
		return mm, nil
	}
	return mm, nil
}

func (mm *matchModel) appendProfile(p Profile) {
	i := 0
	for _, v := range mm.list.Items() {
		if v == nil {
			return
		}
		if v.(Profile).ID == p.ID {
			return
		}
		i++
	}
	mm.list.InsertItem(i, p)
}

func (mm *matchModel) removeProfile(p Profile) {
	for i, v := range mm.list.Items() {
		if v == nil {
			return
		}
		if v.(Profile).ID == p.ID {
			mm.list.RemoveItem(i)
		}
	}
}

func (mm *matchModel) updateProfiles() error {
	ps, err := shellgame.GetMatchingProfiles()
	if err != nil {
		return err
	}
	var profiles []list.Item
	for _, v := range ps {
		profiles = append(profiles, Profile(*v))
	}
	mm.list.SetItems(profiles)
	return nil
}

func (mm *matchModel) createConn() error {
	conn, err := shellgame.ConnectMatchingRoom()
	if err != nil {
		return err
	}
	mm.conn = conn
	return nil
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
			if err := shellgame.WriteConn(mm.conn, m); err != nil {
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
	for {
		msg := &MatchingMsg{}
		if err := shellgame.ReadConn(mm.conn, msg); err != nil {
			return
		}
		p.Send(*msg)
	}
}
