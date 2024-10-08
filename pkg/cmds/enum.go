package cmds

import (
	"github.com/abiosoft/ishell"
)

func Enum(c *ishell.Context) {
	c.Println("> You are likely to be eaten by a grue.")
	c.Println("Type: enum help")
}

func executeWithProgress(script string, c *ishell.Context) {
	c.ProgressBar().Start()
	out, err := execute(script)
	if err != nil {
		c.ProgressBar().Stop()
		c.Println(err.Error())
		return
	}
	c.ProgressBar().Stop()
	c.Println(string(out))
}
