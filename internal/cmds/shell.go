package cmds

import (
	"github.com/abiosoft/ishell"
	"os/exec"
	// TODO make this work again through iShell, for now just execute cmds
	// "github.com/audibleblink/gorsh/internal/shell"
)

func Shell(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println("Usage: shell <cmd> [args]")
		return
	}

	var cmd *exec.Cmd
	if len(c.Args) == 1 {
		cmd = exec.Command(c.Args[0])
	} else {
		argv := strings.Join(c.Args[1:], " ")
		cmd = exec.Command(c.Args[0], argv)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.Println(err.Error())
		return
	}
	c.Println(string(output))
}
