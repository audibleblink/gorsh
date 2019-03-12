package cmds

import (
	"fmt"

	"github.com/audibleblink/HoleySocks/pkg/holeysocks"
	"github.com/audibleblink/gorsh/internal/sshocks"

	"github.com/abiosoft/ishell"
)

func Socks(c *ishell.Context) {
	sshocks.Config.Socks.Remote = fmt.Sprintf("127.0.0.1:%s", c.Args[0])
	err := holeysocks.DarnSocks(sshocks.Config)
	if err != nil {
		c.Println(err.Error())
		return
	}
}
