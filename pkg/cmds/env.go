package cmds

import (
	"github.com/abiosoft/ishell"
	"git.hyrule.link/blink/gorsh/pkg/sitrep"
	"os"
	"strings"
)

func Env(c *ishell.Context) {
	if len(c.Args) > 0 {
		assignment := strings.Join(c.Args, " ")
		input := strings.Split(assignment, "=")
		if len(input) == 2 {
			os.Setenv(input[0], input[1])
			c.Printf("Set: %s=%s\n", input[0], input[1])
		} else {
			c.Println("Usage: env VARNAME=\"some data\"")
		}
		return
	}
	c.Println(sitrep.Environ())
}
