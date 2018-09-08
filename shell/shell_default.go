// +build linux darwin freebsd !windows

package shell

import "os/exec"

func GetShell() *exec.Cmd {
	cmd := exec.Command("/bin/sh")
	return cmd
}
