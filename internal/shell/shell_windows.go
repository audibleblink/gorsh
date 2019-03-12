// +build windows !linux !darwin !freebsd

package shell

import (
	"os/exec"
	"syscall"
)

// TODO make this work again through iShell, for now just execute cmds
func GetShell() *exec.Cmd {
	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
