//go:build windows

package procdump

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var comsvcs = windows.NewLazySystemDLL("comsvcs.dll")
var miniDumpW = comsvcs.NewProc("MiniDumpW")

func Procdump(pid int) error {
	outfile := fmt.Sprintf("%d.mini", pid)
	args := unsafe.Pointer(
		windows.StringToUTF16Ptr(fmt.Sprintf("%d %s full", pid, outfile)),
	)
	_, _, err := miniDumpW.Call(uintptr(0), uintptr(0), uintptr(args))
	if err != syscall.Errno(0) {
		err = fmt.Errorf("procdump failed")
		return err
	}
	return err
}
