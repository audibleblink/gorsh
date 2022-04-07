package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"go-tty"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"github.com/wricardo/gomux"
)

var opts struct {
	Iface  string `short:"i" long:"host" description:"Interface address on which to bind" default:"127.0.0.1" required:"true"`
	Port   string `short:"p" long:"port" description:"Port on which to bind" default:"8443" required:"true"`
	Keys   string `short:"k" long:"keys" description:"Path to folder with server.{pem,key}" default:"./certs" required:"true"`
	Socket string `short:"s" long:"socket" description:"Domain socket from which the program reads"`
}

var sessions = make(map[string]*gomux.Session)

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

	// Since this binary only builds tmux commands and echoes them,
	// it needs to be piped to bash in order to work.
	// Because of this, all logging is sent to stderr
	log.SetOutput(os.Stderr)
}

func main() {
	var listener net.Listener
	var err error

	if opts.Socket == "" {
		// Shell-catching mode. TLS -> TMUX -> Shell
		// Once the shell is caught over TLS, it's unwrapped and sent
		// to a local socket, where it will later be read by a new instance
		// of the server configured to read that socket from within a tmux pane

		listener, err = newTLSListener()
		if err != nil {
			log.Fatal(err)
		}
		log.WithFields(log.Fields{"port": opts.Port, "host": opts.Iface}).Info("Listener started")

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Error(err)
				continue
			}
			defer conn.Close()

			sockF, err := routeToTmux(conn)
			if err != nil {
				log.Error(err)
				continue
			}
			time.Sleep(1 * time.Second) // Give socket time to establish
			go proxyConnToSocket(conn, sockF)
		}

	} else {

		// Post-tmux routing.
		// Creates a socket file and listens for input.
		// If in this branch, binary was started from within tmux.
		// Once the tcp and sockets are mutually proxied with
		// `proxyConnToSocket`, the shell will start
		listener, err = net.Listen("unix", opts.Socket)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()
		log.WithField("socket", opts.Socket).Info("Listener started")

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Error(err)
				continue
			}
			defer conn.Close()
			go startShell(conn)
		}
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

func startShell(conn net.Conn) {
	log.WithFields(log.Fields{"port": opts.Port, "host": opts.Iface}).Info("Incoming")
	reader, writer := bufio.NewReader(conn), bufio.NewWriter(conn)

	// continuously read the incoming data and print to stdout
	go func() {
		for {
			out, err := reader.ReadByte()
			if err != nil {
				log.Fatalf("error reading from implant: %w\n", err)
			}
			fmt.Printf(string(out))
		}
	}()

	teaTeeWhy, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer teaTeeWhy.Close()

	// support resizing
	go func() {
		for ws := range teaTeeWhy.SIGWINCH() {
			fmt.Println("Resized", ws.W, ws.H)
		}
	}()

	restore, err := teaTeeWhy.Raw()
	if err != nil {
		log.Errorf("failed to enter raw: %s\n", err)
		restore()
	}
	defer restore()

	// read stdin from user and send to remote implant
	for {
		key, err := teaTeeWhy.ReadRune()
		if err != nil {
			log.Errorf("failed to read input char: %s\n", err)
			continue
		}

		writer.WriteRune(key)
		writer.Flush()

	}
}

func routeToTmux(conn net.Conn) (string, error) {

	reader := bufio.NewReader(conn)
	hostname, err := reader.ReadString('\n')
	if err != nil {
		err = fmt.Errorf("hostname read failed: %w", err)
		return "", err
	}

	username, err := reader.ReadString('\n')
	if err != nil {
		err = fmt.Errorf("username read failed: %w", err)
		return "", err
	}

	// tmux session names can't contain "."
	hostname = strings.ReplaceAll(strings.TrimSuffix(hostname, "\n"), ".", "_")
	username = strings.TrimSuffix(username, "\n")

	if sessions[hostname] == nil {
		log.WithField("host", hostname).Info("new host connected, creating session")
		sessions[hostname] = gomux.NewSession(hostname, os.Stdout)
	}

	session := sessions[hostname]

	windowName := username
	window := session.AddWindow(windowName)

	if _, err := os.Stat(".state"); os.IsNotExist(err) {
		os.Mkdir(".state", 0700)
	}

	file, err := ioutil.TempFile(".state", fmt.Sprintf("%s.*.sock", windowName))
	if err != nil {
		err = fmt.Errorf("temp file failed: %w", err)
		return "", err
	}
	os.Remove(file.Name())

	path, err := filepath.Abs(file.Name())
	if err != nil {
		err = fmt.Errorf("temp path read failed: %w", err)
		return "", err
	}

	self := os.Args[0]
	cmd := fmt.Sprintf("%s -s %s", self, path)
	log.WithFields(log.Fields{"session": session.Name, "window": windowName}).Info("new shell in tmux")
	window.Exec(`echo -e '\a'`) // ring a bell
	window.Exec(cmd)
	return path, nil
}

func proxyConnToSocket(conn net.Conn, sockF string) {
	socket, err := net.Dial("unix", sockF)
	if err != nil {
		log.WithField("err", err).Error("failed to dial sockF")
		return
	}
	defer socket.Close()
	defer os.Remove(sockF)

	wg := sync.WaitGroup{}

	// forward socket to tcp
	wg.Add(1)
	go (func(socket net.Conn, conn net.Conn) {
		defer conn.Close()
		defer wg.Done()
		io.Copy(conn, socket)
	})(socket, conn)

	// forward tcp to socket
	wg.Add(1)
	go (func(socket net.Conn, conn net.Conn) {
		defer socket.Close()
		defer wg.Done()
		io.Copy(socket, conn)
	})(socket, conn)

	// keep from returning until sockets close so we
	// can cleanup the socket file using `defer`
	wg.Wait()
}

func initSessions() {
	for _, name := range sessionList() {
		sessions[name] = gomux.NewSession(name, os.Stdout)
	}
}

func sessionList() []string {
	// cmd := "tmux list-sessions -F #S"
	cmd := exec.Command("tmux", "list-sessions", "-F", "#S")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic(err)
	}

	return strings.Split(string(output), "\n")
}
