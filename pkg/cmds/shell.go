package cmds

import (
	"git.hyrule.link/blink/gorsh/pkg/myconn"
	"git.hyrule.link/blink/gorsh/pkg/shell"
	"github.com/abiosoft/ishell"
)

func Shell(c *ishell.Context) {
	cmd := shell.GetShell()

	cmd.Stderr = myconn.Conn
	cmd.Stdin = myconn.Conn
	cmd.Stdout = myconn.Conn

	cmd.Run()
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
