package cmds

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"
	"git.hyrule.link/blink/gorsh/pkg/myconn"
	"git.hyrule.link/blink/gorsh/pkg/pivot"
)

func Pivot(c *ishell.Context) {
	connSlice := strings.Split(myconn.ConnectionString, ":")
	// ligolo's default listen port
	host := fmt.Sprintf("%s:11601", connSlice[0])
	if len(c.Args) == 1 {
		host = c.Args[0]
	}
	go pivot.Connect(host)
}
