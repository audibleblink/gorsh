package cmds

import (
	"github.com/abiosoft/ishell"
)

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
		Name:      "mv",
		Help:      "mv <source> <dest>",
		LongHelp:  "Move / rename files",
		Func:      Move,
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
		Aliases:  []string{"sh"},
		Help:     "shell",
		LongHelp: "Drop down to an interactive shell",
		Func:     Shell,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "exec",
		Help:     "exec <args>",
		LongHelp: "Executes a oneoff command",
		Func:     Exec,
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
		Name:     "scan",
		Help:     "scan < ip|cidr > [ 250|500 ]",
		LongHelp: "Ping sweep and TCP scan an IP or subnet. Live hosts only. Defaults to 250",
		Func:     Scanner,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "pivot",
		Help:     "pivot [ip:port]",
		LongHelp: "Connect to a Liglo server to allow for proxying",
		Func:     Pivot,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "find",
		Help:     "find <root> <pattern>",
		LongHelp: "Look for a file name by regex pattern",
		Func:     Find,
	})
}
