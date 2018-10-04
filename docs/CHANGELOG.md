Changes after fork:

* Added a reverse SOCKS5 proxy over ssh. Configure in `configs/ssh.json`
* Uses tmux as a pseudo-C2-like interface, creating a new window with each agent callback
* zipcat: zip > base64 > cat for data small/medium data exfil (zstd/x64 or gzip/x86)
* Remove the use of passing powershell/cmd commands at the main prompt
* Add common file operation commands that use go instead of powershell/cmd
* Download files to victim with HTTP on all platforms
* Download files to victim with SMB on windows agents
* Situational Awareness output on new shells
* Removed Meterpreter functionality
* Removed Shellcode execution
