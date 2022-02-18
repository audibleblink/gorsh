package cmds

import (
	"github.com/abiosoft/ishell"
	"os"
)

func Pwd(c *ishell.Context) {
	output, err := os.Getwd()
	if err != nil {
		c.Println(err.Error())
		return
	}
	c.Println(output)
}
