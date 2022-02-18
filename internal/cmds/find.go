package cmds

import (
	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/find"
)

func Find(c *ishell.Context) {
	if len(c.Args) != 2 {
		c.Println(c.Cmd.Help)
		return
	}

	files, err := find.Find(c.Args[0], c.Args[1])
	if err != nil {
		c.Printf("find failed with %s\n", err)
		return
	}

	for _, file := range files {
		c.Println(file)
	}
}
