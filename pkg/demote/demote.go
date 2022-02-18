//go:build windows

package demote

import (
	"github.com/audibleblink/getsystem"
)

func Demote(pid int) error {
	return getsystem.DemoteProcess(pid)
}
