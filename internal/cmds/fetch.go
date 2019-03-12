package cmds

import (
	"github.com/audibleblink/gorsh/internal/fetch"

	"github.com/abiosoft/ishell"
)

func Fetch(c *ishell.Context) {
	from, to := c.Args[0], c.Args[1]
	bytes, err := fetch.Get(from, to)
	if err != nil {
		c.Println(err.Error())
		return
	}
	c.Printf("Copied %d bytes from %s to %s\n", bytes, from, to)
}
