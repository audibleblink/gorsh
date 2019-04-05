package cmds

import (
	"github.com/abiosoft/ishell"
)

func Enum(c *ishell.Context) {
	c.Println("> You are likely to be eaten by a grue.")
	c.Println("Type: enum help")
}

func executeWithProgress(scriptB64 string, c *ishell.Context) {
	c.ProgressBar().Start()
	out, err := execute(scriptB64)
	if err != nil {
		c.ProgressBar().Stop()
		c.Println(err.Error())
		return
	}
	c.ProgressBar().Stop()
	c.Println(string(out))
}

//stubs for GOOS without a need for these functions
func addSubEnumCmds(sh *ishell.Cmd) *ishell.Cmd {
	return sh
}

//stubs for GOOS without a need for these functions
func execute(scriptB64 string) ([]byte, error) {
	return nil, nil
}
