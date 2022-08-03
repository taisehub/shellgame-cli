package repository

import (
	"net"
)

type ConsoleRepository interface {
	StartShell() (net.Conn, error)
}