// +build linux darwin freebsd !windows

package shell

import "os/exec"

// TODO make this work again through iShell, for now just execute cmds
func GetShell() *exec.Cmd {
	cmd := exec.Command("/bin/sh")
	return cmd
}
