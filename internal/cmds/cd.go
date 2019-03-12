package cmds

import (
	"github.com/abiosoft/ishell"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func Cd(c *ishell.Context) {
	argv := strings.Join(c.Args, " ")
	if argv != "" {
		paths, err := filepath.Glob(argv)
		if err != nil {
			c.Println(err.Error())
			return
		}

		if len(paths) == 1 {
			os.Chdir(paths[0])
		} else {
			c.Println("Glob returned more than 1 result")
			return
		}
	} else {
		usr, _ := user.Current()
		os.Chdir(usr.HomeDir)
	}
	dir, _ := os.Getwd()
	c.Printf("Current Directory: %s\n\n", dir)
}

func CompCd(s []string) []string {
	nodes, _ := ioutil.ReadDir(".")
	var dirs []string
	for _, node := range nodes {
		if node.IsDir() {
			dirs = append(dirs, node.Name())
		}
	}
	return dirs
}
