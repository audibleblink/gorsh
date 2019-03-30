package cmds

import (
	"fmt"
	"os/exec"

	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/enum"
)

func Enum(c *ishell.Context) {
	c.Println("You are likely to be eaten by a grue.")
}

func addSubEnumCmds(sh *ishell.Cmd) {
	sh.AddCmd(&ishell.Cmd{
		Name: "linenum",
		Help: "github.com/rebootuser/linenum",
		Func: func(c *ishell.Context) {
			b64 := enum.LinEnum().Base64()
			executeWithProgress(b64, c)
		},
	})
}

func execute(scriptB64 string) ([]byte, error) {
	cmd := fmt.Sprintf(" echo %s | base64 -d | bash", scriptB64)
	return exec.Command("bash", "-c", cmd).Output()
}

func executeWithProgress(scriptB64 string, c *ishell.Context) {
	c.ProgressBar().Start()
	out, err := execute(scriptB64)
	if err != nil {
		c.ProgressBar().Stop()
		c.Println(err.Error())
		return
	}
	c.ProgressBar().Stop()
	c.Println(out)
}
