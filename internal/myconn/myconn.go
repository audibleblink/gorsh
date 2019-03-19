package myconn

var Conn Writer

type Writer interface {
	Write(s []byte) (int, error)
	Read(s []byte) (int, error)
	Close() error
}

func Send(conn Writer, msg string) {
	conn.Write([]byte(msg))
	conn.Write([]byte("\n"))
}
