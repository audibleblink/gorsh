// +build windows !linux !darwin !freebsd

package shell

import (
	"net"
	"os/exec"
	"syscall"
)

func GetShell() *exec.Cmd {
	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
