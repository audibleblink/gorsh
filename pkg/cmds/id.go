//go:build !windows

package cmds

import (
	"github.com/abiosoft/ishell"
	"git.hyrule.link/blink/gorsh/pkg/sitrep"
)

func Id(c *ishell.Context) {
	output, err := sitrep.UserInfo()
	if err != nil {
		c.Println(err)
		return
	}

	c.Println(output.String())
}
