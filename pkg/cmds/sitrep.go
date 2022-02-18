package cmds

import (
	"github.com/abiosoft/ishell"

	"git.hyrule.link/blink/gorsh/pkg/sitrep"
)

func Sitrep(c *ishell.Context) {
	output := sitrep.SysInfo()
	c.Println(output)
}
