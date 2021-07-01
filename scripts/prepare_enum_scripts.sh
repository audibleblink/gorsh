#!/bin/bash

set -x

powerup=https://raw.githubusercontent.com/PowerShellMafia/PowerSploit/d943001a7defb5e0d1657085a77a0e78609be58f/Privesc/PowerUp.ps1
sherlock=https://raw.githubusercontent.com/rasta-mouse/Sherlock/9f5be56ea2989c01e6ccf19de6b70e0fcd30a11c/Sherlock.ps1
linenum=https://raw.githubusercontent.com/rebootuser/LinEnum/65475312171107e9373dd8b06c9757610f0653d8/LinEnum.sh
winpeas=https://raw.githubusercontent.com/carlospolop/privilege-escalation-awesome-scripts-suite/cc00bf89ab25fc7818aac2a3476539f24c26a720/winPEAS/winPEASbat/winPEAS.bat
linpeas=https://raw.githubusercontent.com/carlospolop/privilege-escalation-awesome-scripts-suite/3e51fee9ac473c70e629ed0c3a0f79386165aa9f/linPEAS/linpeas.sh
jaws=https://raw.githubusercontent.com/411Hall/JAWS/233f142fcb1488172aa74228a666f6b3c5c48f1d/jaws-enum.ps1

dest=internal/enum/embed

function download() {
  wget "${1}" -O "${dest}/${2}"
}

## PowerUp
download $powerup powerup.ps1 && \
echo -e \\nInvoke-AllChecks >> ${dest}/powerup.ps1

## WinPEAS
download $winpeas winpeas.bat

## Jaws
download $jaws jaws.ps1

## LinPEAS
download $linpeas linpeas.sh

## Sherlock
download $sherlock sherlock.ps1 && \
echo -e \\nFind-AllVulns >> ${dest}/sherlock.ps1

## LinEnum
download $linenum linenum.sh
sed -i '2ithorough=1' ${dest}/linenum.sh
