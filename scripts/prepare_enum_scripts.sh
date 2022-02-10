#!/bin/bash

set -x

linenum=https://raw.githubusercontent.com/rebootuser/LinEnum/65475312171107e9373dd8b06c9757610f0653d8/LinEnum.sh
linpeas=https://raw.githubusercontent.com/carlospolop/privilege-escalation-awesome-scripts-suite/3e51fee9ac473c70e629ed0c3a0f79386165aa9f/linPEAS/linpeas.sh

dest=internal/enum/embed

function download() {
  wget "${1}" -O "${dest}/${2}"
}

## LinPEAS
download $linpeas linpeas.sh

## LinEnum
download $linenum linenum.sh
sed -i '2ithorough=1' ${dest}/linenum.sh
