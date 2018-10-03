package socks

import (
	"fmt"

	"github.com/armon/go-socks5"
	"github.com/audibleblink/gorsh/internal/sshtacr"
)

func ListenAndForward(port string) {
	// Create a SOCKS5 server
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost
	connectString := fmt.Sprintf("127.0.0.1:%s", port)
	go server.ListenAndServe("tcp", connectString)
	go sshtacr.ForwardService(port)
}
