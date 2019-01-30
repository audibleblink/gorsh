package sshtacr

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/gobuffalo/packr"
	"golang.org/x/crypto/ssh"
)

type sshServer struct {
	Username string `json:"username"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func (s *sshServer) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
		}
		chDone <- true
	}()

	<-chDone
}

// ForwardService implements reverse port forwarding similar to the -R flag
// in openssh-client. Configuration is done in the /configs/ssh.json file.
// NOTE The generated keys and ssh.json data are embedded in the binary so
// take the appropriate precautions when setting up the ssh server user.
func ForwardService(port string) error {

	// unpack the configs and ssh keys from the binary
	// that were packed at compile-time
	box := packr.NewBox("../../configs")
	privateKeyString, err := box.Find("id_ed25519")
	configs, err := box.Find("ssh.json")
	if err != nil {
		return err
	}

	server := sshServer{}
	if err := json.Unmarshal(configs, &server); err != nil {
		return err
	}

	privateKey, err := ssh.ParsePrivateKey(privateKeyString)
	auth := ssh.PublicKeys(privateKey)

	sshConfig := &ssh.ClientConfig{
		User:            server.Username,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH server
	serverConn, err := ssh.Dial("tcp", server.String(), sshConfig)
	if err != nil {
		return fmt.Errorf("Dial INTO remote server error: %s", err)
	}

	// Publish the designated local port to the same port on the remote SSH server
	connectStr := fmt.Sprintf("127.0.0.1:%s", port)
	listener, err := serverConn.Listen("tcp", connectStr)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("INFO: %s", err))
	}
	defer listener.Close()
	return err
}
