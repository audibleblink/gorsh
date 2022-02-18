package cmds

import (
	"fmt"
	"os/exec"

	"github.com/abiosoft/ishell"
	"git.hyrule.link/blink/gorsh/pkg/enum"
)

func addSubEnumCmds(sh *ishell.Cmd) *ishell.Cmd {
	sh.AddCmd(&ishell.Cmd{
		Name: "linenum",
		Help: "github.com/rebootuser/linenum",
		Func: func(c *ishell.Context) {
			b64 := enum.LinEnum().Base64()
			executeWithProgress(b64, c)
		},
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "linpeas",
		Help: "github.com/carlospolop",
		Func: func(c *ishell.Context) {
			b64 := enum.LinPeas().Base64()
			executeWithProgress(b64, c)
		},
	})
	return sh
}

func execute(scriptB64 string) ([]byte, error) {
	cmd := fmt.Sprintf(" echo %s | base64 -d | bash", scriptB64)
	return exec.Command("bash", "-c", cmd).Output()
}
