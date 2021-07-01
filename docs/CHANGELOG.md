Changes after fork:

[01 Jul 21]
* extract connection logic to `core` module
* holeysocks case change and update 
* semantic use of make targets and dependencies
* adds shellcode target and dllmain code to Makefile
* improve handling of cgo when opting in to zstd
* replace vsfgen with native go embed (go1.16+)

[18 Mar 19]
* adds enumeration module with OS-specific scripts
    - Linux: LinEnum by @rebootuser
    - Windows: PowerUp, JAWS, Sherlock

[15 Mar 19]

* `spawn` command creates a new shell for redundancy
* `spawn host:port` creats a new shell to different hosts
* agent re-writen to use ishell for easier command additions and tab-completion
* addition of tty-capable lister that provides emacs movements and other readline capabilities
* vi-mode cli editing
* tab-completion
* reworks the tmux workflow to use the new gorsh-listener
* use unix sockets to receive plaintext comms from agents when running behind a reverse proxy
* `env` command expanded to also set variables
* `cp` added
* converted `shell` to a non-interactive,  one-off code executor


[older]
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
