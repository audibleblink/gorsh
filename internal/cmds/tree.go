package cmds

import (
	"fmt"
	"github.com/abiosoft/ishell"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Tree(c *ishell.Context) {
	argv := strings.Join(c.Args, " ")
	path := argv
	if argv == "" {
		path = "."
	}

	c.Println(tree(path, "", ""))
}

func tree(root, indent, result string) string {
	fi, err := os.Stat(root)
	if err != nil {
		return fmt.Sprintf("could not stat %s: %v\n", root, err.Error())
	}

	if !fi.IsDir() {
		return ""
	}

	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Sprintf("could not read dir %s: %v\n", root, err.Error())
	}

	var names []string
	for _, fi := range fis {
		names = append(names, fi.Name())
	}

	var out string
	for i, name := range names {
		var add string
		if i == len(names)-1 {
			out += fmt.Sprintf("%s+-- %s\n", indent, name)
			//BUG TODO these don't display well when sent over tcp
			// out += fmt.Sprintf("%s└── %s\n", indent, name
			add = "    "
		} else {
			//BUG TODO these don't display well when sent over tcp
			// out += fmt.Sprintf("%s├── %s\n", indent, name)
			// add = "│   "
			out += fmt.Sprintf("%s|-- %s\n", indent, name)
			add = "|   "
		}

		out += tree(filepath.Join(root, name), indent+add, result)
	}
	return out
}
