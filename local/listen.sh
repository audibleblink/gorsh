#!/bin/bash

SOCKDIR="$(mktemp -d)"
SOCKF="${SOCKDIR}/usock"

# Start tmux, if needed
# tmux new -s GOSH
tmux has-session -t GOSH || tmux new-session -d -s GOSH

IP=$(lsof -Pni | grep "socat.*$PORT" | tail -n 1 | sed 's/>/ /g' | awk '{ print $10 }')
# Create window
tmux new-window -n "$IP" "socat UNIX-LISTEN:${SOCKF},umask=0077 STDIO"

# Wait for socket
while test ! -e ${SOCKF} ; do sleep 1 ; done

# Use socat to ship data between the unix socket and STDIO.
exec socat STDIO UNIX-CONNECT:${SOCKF}
