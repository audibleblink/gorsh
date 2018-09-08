#!/bin/bash

SOCKDIR="$(mktemp -d)"
SOCKF="${SOCKDIR}/usock"

# Start tmux, if needed
# tmux new -s GOSH
tmux has-session -t GOSH || tmux new-session -d -s GOSH

# Create window
tmux new-window "socat UNIX-LISTEN:${SOCKF},umask=0077 STDIO"

# Wait for socket
while test ! -e ${SOCKF} ; do sleep 1 ; done

# Use socat to ship data between the unix socket and STDIO.
exec socat STDIO UNIX-CONNECT:${SOCKF}
