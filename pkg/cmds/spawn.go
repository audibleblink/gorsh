package cmds

import (
	"os"
	"os/exec"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/shirou/gopsutil/process"
)

func Spawn(c *ishell.Context) {
	pid := int32(os.Getpid())
	proc, err := process.NewProcess(pid)
	if err != nil {
		c.Println(err.Error())
		return
	}

	p, err := proc.Exe()
	if err != nil {
		c.Println(err.Error())
		return
	}

	var cmd *exec.Cmd
	if len(c.Args) > 0 {
		argv := strings.Join(c.Args, " ")
		cmd = exec.Command(p, "--connect", argv)
	} else {
		cmd = exec.Command(p)
	}

	err = cmd.Start()
	if err != nil {
		c.Println(err.Error())
	} else {
		c.Println("Spawn command completed")
	}
}
