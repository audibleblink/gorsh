//go:build !windows

package shell

import (
	"os/exec"
	"syscall"
)

func GetShell() *exec.Cmd {
	cmd := exec.Command("/bin/sh")
	return cmd
}

func BGExec(prog string, args []string) (int, error) {
	return syscall.ForkExec(prog, args, nil)
}
