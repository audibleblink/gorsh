package cmds

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/myconn"
	"github.com/audibleblink/gorsh/internal/pivot"
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
