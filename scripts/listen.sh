#!/bin/bash -xe

NAME=GORSH
SOCKDIR="$(mktemp -d)"
SOCKF="${SOCKDIR}/usock"

function finish () {
        rm -rf $SOCKDIR
}
trap finish EXIT

tmux set-hook -g session-created 'set remain-on-exit on'

# Start tmux in the background, if needed
tmux has-session -t "${NAME}" || tmux new-session -d -s "${NAME}"

# Start the listener in a tmux pane, using the new Unix domain socket, and wait
tmux new-window -t "${NAME}" -a -n "shell" "stty -icanon -brkint -isig && build/srv/gorsh-server -s ${SOCKF}"
sleep 1

# Hook up the STDIN from the calling socat instance to the Unix socket that gorsh is listening to
socat - UNIX:"${SOCKF}"
