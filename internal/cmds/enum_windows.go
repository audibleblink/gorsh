package cmds

import (
	"github.com/abiosoft/ishell"
	"github.com/audibleblink/gorsh/internal/enum"
)

func addSubEnumCmds(sh *ishell.Cmd) *ishell.Cmd {
	return sh
}

func enumFn(fn func() *enum.EnumScript) func(*ishell.Context) {
	return func(c *ishell.Context) {}
}

func execute(b64Script string) ([]byte, error) {
	return []byte{}, nil
}
