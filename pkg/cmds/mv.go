package cmds

import (
	"os"

	"github.com/abiosoft/ishell"
)

func Move(c *ishell.Context) {
	src, dst := c.Args[0], c.Args[1]
	err := os.Rename(src, dst)
	if err != nil {
		c.Printf("move failed: %s\n", err)
		return
	}
}
