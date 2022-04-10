package cmds

import (
	"os/exec"
	"strings"

	"git.hyrule.link/blink/gorsh/pkg/enum"
	"github.com/abiosoft/ishell"
)

func addSubEnumCmds(sh *ishell.Cmd) *ishell.Cmd {
	sh.AddCmd(&ishell.Cmd{
		Name: "linenum",
		Help: "github.com/rebootuser/linenum",
		Func: func(c *ishell.Context) {
			scrpt := enum.LinEnum().String()
			executeWithProgress(scrpt, c)
		},
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "linpeas",
		Help: "github.com/carlospolop",
		Func: func(c *ishell.Context) {
			scrpt := enum.LinPeas().String()
			executeWithProgress(scrpt, c)
		},
	})
	return sh
}

func execute(script string) ([]byte, error) {
	c := exec.Command("/bin/bash")
	c.Stdin = strings.NewReader(script)
	return c.Output()
}
