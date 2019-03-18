package cmds

import (
	"fmt"
	"os/exec"

	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/enum"
)

func Enum(c *ishell.Context) {
	var script string
	var err error

	choice := c.MultiChoice([]string{
		"linenum",
	}, "Run which script?")

	switch choice {
	case 0:
		script = enum.LinEnum().Base64()
	default:
		c.Println("You are likely to be eaten by a grue")
		return
	}

	if err != nil {
		c.Println(err.Error())
		return
	}

	c.ProgressBar().Start()
	out, err := execute(script)
	if err != nil {
		c.ProgressBar().Stop()
		c.Println(err.Error())
		return
	}
	c.ProgressBar().Stop()
	c.Println(string(out))
}

func execute(scriptB64 string) ([]byte, error) {
	cmd := fmt.Sprintf(" echo %s | base64 -d | bash", scriptB64)
	return exec.Command("bash", "-c", cmd).Output()
}
