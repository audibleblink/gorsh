package cmds

import (
	"fmt"
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

	argv := strings.Join(c.Args, " ")
	connInfo := strings.Split(argv, ":")

	spawnArgs := fmt.Sprintf("--host %s --port %s", connInfo[0], connInfo[1])
	cmd := exec.Command(p, spawnArgs)
	cmd.Start()
}
