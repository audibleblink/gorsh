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

func ForwardService(port string) {

	connectStr := fmt.Sprintf("127.0.0.1:%s", port)
	box := packr.NewBox("../../configs")
	privateKeyString := box.Bytes("id_ed25519")
	configs := box.Bytes("ssh.json")
	server := sshServer{}
	if err := json.Unmarshal(configs, &server); err != nil {
		panic(err)
	}

	privateKey, err := ssh.ParsePrivateKey(privateKeyString)
	auth := ssh.PublicKeys(privateKey)

	sshConfig := &ssh.ClientConfig{
		User:            server.Username,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH remote server using serverEndpoint
	serverConn, err := ssh.Dial("tcp", server.String(), sshConfig)
	if err != nil {
		log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
	}

	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", connectStr)
	if err != nil {
		log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
	}
	defer listener.Close()

	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, err := net.Dial("tcp", connectStr)
		if err != nil {
			log.Fatalln(fmt.Printf("Dial INTO local service error: %s", err))
		}

		client, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		handleClient(client, local)
	}

}
