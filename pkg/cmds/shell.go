package cmds

import (
	"git.i.ctrl.red/blink/gorsh/pkg/shell"
	"github.com/abiosoft/ishell"
)

func Shell(c *ishell.Context) {
	sh, err := shell.GetShell()
	if err != nil {
		c.Print(err)
		return
	}
	sh.Wait()
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
