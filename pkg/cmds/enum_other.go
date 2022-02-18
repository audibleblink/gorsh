// +build !linux,!windows

package cmds

import (
	"github.com/abiosoft/ishell"
)

//stubs for GOOS without a need for these functions
func addSubEnumCmds(sh *ishell.Cmd) *ishell.Cmd {
	return sh
}

//stubs for GOOS without a need for these functions
func execute(scriptB64 string) ([]byte, error) {
	return nil, nil
}
