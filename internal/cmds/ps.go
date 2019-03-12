package cmds

import (
	"github.com/abiosoft/ishell"

	"github.com/audibleblink/gorsh/internal/sitrep"
)

func Ps(c *ishell.Context) {
	output := sitrep.Processes()
	c.Println(output)
}
