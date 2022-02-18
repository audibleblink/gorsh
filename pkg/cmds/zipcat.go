package cmds

import (
	"encoding/base64"
	"strings"

	"github.com/abiosoft/ishell"

	"git.hyrule.link/blink/gorsh/pkg/zip"
)

func Zipcat(c *ishell.Context) {
	argv := strings.Join(c.Args, " ")
	zipped, err := zip.ZipWriter(argv)
	if err != nil {
		c.Println(err.Error())
		return
	}

	b64 := base64.StdEncoding.EncodeToString(zipped)
	// b64 := base64.StdEncoding.EncodeToString(zipped.Read())
	c.Println(b64)
	c.Println("")
}
