#!/bin/sh

socat OPENSSL-LISTEN:$PORT,fork,key=$KEY,cert=$PEM,reuseaddr,verify=0 EXEC:$LISTEN
