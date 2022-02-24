package myconn

import (
	"io"
	"strings"
)

var (
	Conn             io.ReadWriteCloser
	ConnectionString string
)

type Writer interface {
	Write(s []byte) (int, error)
	Read(s []byte) (int, error)
	Close() error
}

func Send(conn Writer, msg string) {
	conn.Write([]byte(msg))
	conn.Write([]byte("\n"))
}

func Host() string {
	split := strings.Split(ConnectionString, ":")
	return split[0]
}
