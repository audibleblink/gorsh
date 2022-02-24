package cmds

import (
	"git.hyrule.link/blink/gorsh/pkg/fetch"

	"github.com/abiosoft/ishell"
)

func Fetch(c *ishell.Context) {
	if len(c.Args) != 2 {
		c.Println(c.Cmd.LongHelp)
		return
	}
	from, to := c.Args[0], c.Args[1]
	bytes, err := fetch.Get(from, to)
	if err != nil {
		c.Println(err)
		return
	}
	c.Printf("Copied %d bytes from %s to %s\n", bytes, from, to)
}
