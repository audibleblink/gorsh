package cmds

import (
	"git.hyrule.link/blink/gorsh/pkg/sitrep"
	"github.com/abiosoft/ishell"
	"github.com/audibleblink/getsystem"
	"golang.org/x/sys/windows"
)

func Id(c *ishell.Context) {
	output, err := sitrep.UserInfo()
	if err != nil {
		c.Println(err)
		return
	}

	procH := windows.GetCurrentProcessToken()
	output.Token, err = getsystem.TokenOwner(procH)
	if err != nil {
		c.Println(err)
		return
	}

	c.Println(output.String())
}
