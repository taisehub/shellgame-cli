package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"io"
)

type profile struct {
	ID   string
	Name string
}

func (p profile) FilterValue() string { return "" }

type profileDelegate struct{}

func (d profileDelegate) Height() int                               { return 1 }
func (d profileDelegate) Spacing() int                              { return 0 }
func (d profileDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d profileDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(profile)
	if !ok {
		return
	}

	str := fmt.Sprintf("* %s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render(">  " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

// マッチ画面を表すtea.Modelの実装
// tea.Modelの実装にしないとlist.ModelのIndex()メソッドが動作しない
type matchModel struct {
	parent  *model
	waiting bool
	list    list.Model
}

func NewMatchModel() matchModel {
	profiles := []list.Item{
		//後で消す。
		profile{ID: "11111", Name: "test1"},
		profile{ID: "22222", Name: "test2"},
		profile{ID: "33333", Name: "test3"},
	}
	l := list.New(profiles, profileDelegate{}, width, 14)
	l.Title = "対戦相手を選択してください"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return matchModel{list: l, waiting: false}
}

func (mm matchModel) Init() tea.Cmd {
	// websocketを用いて接続するクライアントを起動する。
	// websocketクライアントは、データの送信と受信するメソッドを持ちそれぞれGoroutineで管理する。
	//　対戦待ちのプレイヤーを取得するHTTPリクエストを送信して、リストに反映する。
	// mm.list.SetItems([]profile...)
	return nil
}

func (mm *matchModel) setParent(p *model) {
	mm.parent = p
}

func (mm matchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//case websocketMsg: // websocketからjsonを受け取る
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
			// 選択したプロフィールの取り方
			// profile, _ := mm.list.SelectedItem().(profile)
			//
			// シェルゲーサーバに選択した対戦相手を含む対戦要求データを送信。
			// 送信時に3分後にtimeoutMsgを通知する処理をgoroutineで動かす。
			// 送信するデータの例
			//{
			//		"source": { //自分のProfile
			//   		"id": "1-1-1",
			//   		"name": "player1",
			//  	},
			//  	"dest": { //対戦相手のProfile
			//  		"id": "2-2-2",
			//  		"name": "player2",
			//  	}
			//  	"data": 1,
			//}

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
