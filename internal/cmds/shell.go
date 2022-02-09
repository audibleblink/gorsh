package cmds

import (
	"os"

	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/shell"
)

func Shell(c *ishell.Context) {
	cmd := shell.GetShell()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func Exec(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println("Usage: shell <cmd> [args]")
		return
	}

	cmd := shell.GetShell()
	cmd.Args = append(cmd.Args, "-c")
	cmd.Args = append(cmd.Args, c.Args...)
	err := cmd.Start()
	if err != nil {
		c.Println(err.Error())
		return
	}
	c.Printf("started %d\n", cmd.Process.Pid)
}
