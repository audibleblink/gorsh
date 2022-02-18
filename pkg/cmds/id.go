//go:build !windows

package cmds

import (
	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/sitrep"
)

func Id(c *ishell.Context) {
	output, err := sitrep.UserInfo()
	if err != nil {
		c.Println(err)
		return
	}

	c.Println(output.String())
}
