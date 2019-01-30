package socks

import (
	"fmt"

	"github.com/armon/go-socks5"
	"github.com/audibleblink/gorsh/internal/sshtacr"
)

// ListenAndForward creates a SOCKS5 proxy and publishes the port using
// reverse SSH tunnels. The receiving SSH server can then proxy through
// the network by using the now-local SOCKS proxy into the target network
func ListenAndForward(port string) error {
	// Create a SOCKS5 server
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		return err
	}

	// Create SOCKS5 proxy on localhost
	result := make(chan error)
	connectString := fmt.Sprintf("127.0.0.1:%s", port)

	go func() {
		go server.ListenAndServe("tcp", connectString)

		err = sshtacr.ForwardService(port)
		if err != nil {
			result <- err
		}

	}()

	return <-result
}
