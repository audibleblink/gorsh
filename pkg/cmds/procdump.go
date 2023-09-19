//go:build windows

package cmds

import (
	"strconv"

	"git.i.ctrl.red/blink/gorsh/pkg/procdump"
	"github.com/abiosoft/ishell"
)

func Procdump(c *ishell.Context) {
	if len(c.Args) != 1 {
		c.Println(c.Cmd.Help)
		return
	}

	pid, err := strconv.Atoi(c.Args[0])
	if err != nil {
		c.Println("failed to convert to int")
		return
	}

	err = procdump.Procdump(pid)
}
