#!/bin/bash

powerup=https://raw.githubusercontent.com/PowerShellMafia/PowerSploit/863699d97e55fe375fc67ada9e3d99d462cbe1d0/Privesc/PowerUp.ps1
jaws=https://raw.githubusercontent.com/belane/JAWS/88c8f21812de58a04f2a4874ea359356277ac089/jaws-enum.ps1
sherlock=https://raw.githubusercontent.com/rasta-mouse/Sherlock/9f5be56ea2989c01e6ccf19de6b70e0fcd30a11c/Sherlock.ps1
linenum=https://raw.githubusercontent.com/rebootuser/LinEnum/ddfd743d817e3c189da567191d8488ecbfd65f1f/LinEnum.sh

function download() {
  curl -s "${1}" > scripts/"${2}"
}

## PowerUp
download $powerup powerup.ps1 && \
echo -e \\nInvoke-AllChecks >> scripts/powerup.ps1

## JAWS
download $jaws jaws.ps1

## Sherlock
download $sherlock sherlock.ps1 && \
echo -e \\nFind-AllVulns >> scripts/sherlock.ps1

## LinEnum
download $linenum linenum.sh
sed -i '2ithorough=1' scripts/linenum.sh
