package cmds

import (
	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/pivot"
)

func Pivot(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println("requires host:port")
		return
	}

	host := c.Args[0]
	go pivot.Connect(host)

}
