package tokens

import (
	"github.com/audibleblink/getsystem"
	"golang.org/x/sys/windows"
)

func StealToken(pid int) error {
	return getsystem.OnThread(pid)
}

func RevToSelf() {
	windows.RevertToSelf()
}
