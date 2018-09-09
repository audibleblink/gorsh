#!/bin/bash

SOCKDIR="$(mktemp -d)"
SOCKF="${SOCKDIR}/usock"
NAME=GORSH

# Start tmux in the background, if needed
tmux has-session -t "${NAME}" || tmux new-session -d -s "${NAME}"

# Create window in the created sessions
IP=$(lsof -Pni | grep "socat.*$PORT" | tail -n 1 | sed 's/>/ /g' | awk '{ print $10 }')
tmux  new-window -t "${NAME}:" -a -n "$IP" "socat UNIX-LISTEN:${SOCKF},umask=0077 STDIO"

# Wait for socket
while test ! -e ${SOCKF} ; do sleep 1 ; done

# Use socat to ship data between the unix socket and STDIO.
exec socat STDIO UNIX-CONNECT:${SOCKF}
