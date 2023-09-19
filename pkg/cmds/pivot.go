package cmds

import (
	"fmt"

	"git.i.ctrl.red/blink/gorsh/pkg/myconn"
	"git.i.ctrl.red/blink/gorsh/pkg/pivot"
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
