package cmds

import (
	"github.com/abiosoft/ishell"
	"os"
	"strings"
)

func Rm(c *ishell.Context) {
	argv := strings.Join(c.Args, " ")
	err := os.RemoveAll(argv)
	if err != nil {
		c.Println(err.Error())
		return
	}
	c.Printf("%s Deleted.\n", argv)
}
