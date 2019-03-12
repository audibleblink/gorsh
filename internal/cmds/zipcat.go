package cmds

import (
	"encoding/base64"
	"github.com/abiosoft/ishell"
	"strings"

	"github.com/audibleblink/gorsh/internal/zip"
)

func Zipcat(c *ishell.Context) {
	argv := strings.Join(c.Args, " ")
	bytes, err := zip.Bytes(argv)
	if err != nil {
		c.Println(err.Error())
		return
	}
	b64 := base64.StdEncoding.EncodeToString(bytes)
	c.Println(b64)
	c.Println("")
}
