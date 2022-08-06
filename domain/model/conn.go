package model

type Conn interface {
	Close() error
	Write(Message) error
	Read(Message) error
}
