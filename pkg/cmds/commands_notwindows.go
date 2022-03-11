//go:build !windows

package cmds

import (
	"strconv"
	"syscall"

	"github.com/abiosoft/ishell"
)

func RegisterNotWindowsCommands(sh *ishell.Shell) {
	sh.AddCmd(&ishell.Cmd{
		Name:     "setuid",
		LongHelp: "setuid to 0 if bin is suid",
		Help:     "setuid",
		Func: func(c *ishell.Context) {
			id := 0
			if len(c.Args) > 0 {
				id, _ = strconv.Atoi(c.Args[0])
			}
			err := syscall.Setuid(id)
			if err != nil {
				c.Printf("failure: %s\n", err)
				return
			}
		},
	})
}

func RegisterWindowsCommands(sh *ishell.Shell) {}
