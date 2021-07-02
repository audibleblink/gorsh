package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path"

	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-tty"
	log "github.com/sirupsen/logrus"
)

var opts struct {
	Iface  string `short:"i" long:"host" description:"Interface address on which to bind" default:"127.0.0.1" required:"true"`
	Port   string `short:"p" long:"port" description:"Port on which to bind" default:"8443" required:"true"`
	Keys   string `short:"k" long:"keys" description:"Path to folder with server.{pem,key}" default:"./certs" required:"true"`
	Socket string `short:"s" long:"socket" description:"Domain socket from which the program reads"`
}

func init() {
	_, err := flags.Parse(&opts)
	// the flags package returns an error when calling --help for
	// some reason so we look for that and exit gracefully
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}
}

func main() {
	var listener net.Listener
	var err error

	if opts.Socket == "" {
		listener, err = newTLSListener()
		log.WithFields(log.Fields{"port": opts.Port, "host": opts.Iface}).Info("Listener started")
	} else {
		listener, err = net.Listen("unix", opts.Socket)
		log.WithFields(log.Fields{"socket": opts.Socket}).Info("Listener started")
	}
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			log.Error(err)
			continue
		}
		go handleConnection(conn)
	}
}

func newTLSListener() (net.Listener, error) {
	pem := path.Join(opts.Keys, "server.pem")
	key := path.Join(opts.Keys, "server.key")
	cer, err := tls.LoadX509KeyPair(pem, key)
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	connStr := fmt.Sprintf("%s:%s", opts.Iface, opts.Port)
	return tls.Listen("tcp", connStr, config)
}

func handleConnection(conn net.Conn) {
	log.WithFields(log.Fields{"port": opts.Port, "host": opts.Iface}).Info("Incoming")
	reader, writer := bufio.NewReader(conn), bufio.NewWriter(conn)

	// A.B.C - Always Be Checking if there's new data to pull down on the wire
	go func() {
		for {
			out, err := reader.ReadByte()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(string(out))
		}
	}()

	teaTeeWhy, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer teaTeeWhy.Close()

	go func() {
		for ws := range teaTeeWhy.SIGWINCH() {
			fmt.Println("Resized", ws.W, ws.H)
		}
	}()

	for {
		key, err := teaTeeWhy.ReadRune()
		if err != nil {
			log.Fatal(err)
		}

		writer.WriteRune(key)
		writer.Flush()
	}
}
