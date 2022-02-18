package cmds

import (
	"io/ioutil"
	"strings"

	"github.com/abiosoft/ishell"
	"git.hyrule.link/blink/gorsh/pkg/directory"
)

func Ls(c *ishell.Context) {
	argv := strings.Join(c.Args, " ")
	output, err := directory.List(argv)
	if err != nil {
		c.Println(err.Error())
		return
	} else {
		c.Println(output)
	}
}

func CompLs(s []string) []string {
	nodes, _ := ioutil.ReadDir(".")
	var names []string
	for _, node := range nodes {
		name := node.Name()
		if node.IsDir() {
			name = name + "/"
		}
		names = append(names, name)
	}
	return names
}
