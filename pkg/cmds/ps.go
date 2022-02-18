package cmds

import (
	"github.com/abiosoft/ishell"
	"git.hyrule.link/blink/gorsh/pkg/sitrep"
)

func Ps(c *ishell.Context) {
	processes, err := sitrep.ProcessInfo()
	if err != nil {
		c.Printf("could not load processes: %s", err)
		return
	}

	if doVerbose(c) {
		for _, process := range processes {
			c.Println(process.String())
		}
	} else {
		for _, process := range processes {
			c.Println(process.ConciseString())
		}
	}
}

func doVerbose(c *ishell.Context) (v bool) {
	return len(c.Args) > 0
}
