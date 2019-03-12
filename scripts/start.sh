#!/bin/sh

socat -d -d OPENSSL-LISTEN:$PORT,fork,key=$KEY,cert=$PEM,reuseaddr,verify=0 EXEC:$LISTEN
