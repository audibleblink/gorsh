#!/bin/bash

SOCKDIR="$(mktemp -d)"
SOCKF="${SOCKDIR}/usock"
NAME=GORSH

function finish () {
        rm -rf $SOCKDIR
}
trap finish EXIT

function kill_window() {
        tmux send -t $NAME "tmux kill-session -t $1" ENTER
        exit 1
}

# Start tmux in the background, if needed
tmux has-session -t "${NAME}" || tmux new-session -d -s "${NAME}"

# Create window in the created session and start a socket listener connected to stdio
IP=$(lsof -Pni | grep "socat.*$PORT" | tail -n 1 | sed 's/>/ /g' | awk '{ print $10 }')
tmux new-window -t "${NAME}" -a -n "$IP" "socat UNIX-LISTEN:${SOCKF},umask=0077 STDIO"

# Wait 5 seconds for shell to come in; kill listener/window otherwise
breaker=0
while :; do
        [[ -e ${SOCKF} ]] && break
        sleep 1
        breaker=$[breaker+1]
        if [[ $breaker -ge 5 ]]; then
                sess="$(tmux list-windows -t ${NAME} | cut -d ':' -f 1 | tail -1)"
                kill_window "${sess}"
        fi
done

# Hook up stdio to the listening socket in the tmux session
socat STDIO UNIX-CONNECT:${SOCKF}
