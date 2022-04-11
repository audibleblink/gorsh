//go:build !windows

package shell

import (
	"os/exec"
	"syscall"
)

func GetShell() *exec.Cmd {
	// must be pty in order to interact remotely
	// TODO cycle through many types of pty invocations
	// in case script isn't present
	cmd := exec.Command("/usr/bin/script", "-qc", "/bin/bash", "/dev/null")
	return cmd
}

func BGExec(prog string, args []string) (int, error) {
	return syscall.ForkExec(prog, args, nil)
}
