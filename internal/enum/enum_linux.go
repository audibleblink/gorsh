//build !windows !darwin

package enum

import _ "embed"

//go:embed embed/linenum.sh
var linEnum []byte

func LinEnum() *EnumScript {
	return &EnumScript{Data: linEnum}
}
