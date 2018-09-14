Changes after fork:

* Uses tmux as a pseudo-C2-like interface, creating a new window with each agent callback
* zipcat: zip > base64 > cat for data small/medium data exfil (zstd/x64 or gzip/x86)
* Download files to victim with HTTP on all platform
* Download files to victim with SMB on all windows
* Situational Awareness output on new shells
* Removed Meterpreter functionality
* Removed Shellcode execution
* Remove the use of passing power/shell commands at the gorsh prompt
* Add common file operation commands that use go instead of power/shell
