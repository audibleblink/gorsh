#!/bin/sh
ncat --ssl --ssl-cert server.pem --ssl-key server.key -lvp 4445
