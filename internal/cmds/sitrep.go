package cmds

import (
	"github.com/abiosoft/ishell"

	"github.com/audibleblink/gorsh/internal/sitrep"
)

func Sitrep(c *ishell.Context) {
	output := sitrep.SysInfo()
	c.Println(output)
}
