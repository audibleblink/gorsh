#!/bin/bash

SOCKDIR="$(mktemp -d)"
SOCKF="${SOCKDIR}/usock"
NAME=GORSH

# Start tmux in the background, if needed
tmux has-session -t "${NAME}" || tmux new-session -d -s "${NAME}"

# Create window in the created sessions
tmux  new-window -t "${NAME}" "socat UNIX-LISTEN:${SOCKF},umask=0077 STDIO"

# Wait for socket
while test ! -e ${SOCKF} ; do sleep 1 ; done

# Use socat to ship data between the unix socket and STDIO.
exec socat STDIO UNIX-CONNECT:${SOCKF}
