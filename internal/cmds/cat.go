package cmds

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/abiosoft/ishell"
)

type cmdFunc func(...string) string
type genFunc func(string) ([]byte, error)

func Cat(c *ishell.Context) {
	argv := strings.Join(c.Args, " ")
	output, err := handleGlob(argv, ioutil.ReadFile)
	if err != nil {
		c.Println(err.Error())
		return
	}
	c.Println(output)
}

func handleGlob(path string, cb genFunc) (string, error) {
	matches, err := filepath.Glob(path)
	if err != nil {
		return "", err
	}

	var result string
	var errors string
	for _, file := range matches {
		output, err := ioutil.ReadFile(file)
		if err != nil {
			errors += fmt.Sprintf("%s\n", err.Error())
		}

		result += fmt.Sprintf("%s\n", output)
	}
	return result + errors, err
}
