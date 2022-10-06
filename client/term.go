package shellgame

import (
	"golang.org/x/term"
	"io"
	"os"
)

// Terminalは.github.com/charmbracelet/bubbletea.ExecCommandの実装
// websoketを利用してシェルゲーサーバで用意されるコンテナに接続する。
type Terminal struct {
	Stdin  io.Reader
	Stdout io.Writer
}

func (t *Terminal) Run() error {
	wsconn, err := ConnectShell()
	if err != nil {
		return err
	}
	defer wsconn.Close()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	conn := wsconn.UnderlyingConn()
	go func() { io.Copy(conn, t.Stdin) }()
	io.Copy(t.Stdout, conn)
	return nil
}

func (t *Terminal) SetStdin(r io.Reader) {
	if t.Stdin == nil {
		t.Stdin = r
	}
}

func (t *Terminal) SetStdout(w io.Writer) {
	if t.Stdout == nil {
		t.Stdout = w
	}
}

func (t *Terminal) SetStderr(w io.Writer) {
	return
}
