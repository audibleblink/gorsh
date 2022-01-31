//build !windows !darwin

package enum

import _ "embed"

//go:embed embed/linenum.sh
var linEnum []byte

//go:embed embed/linpeas.sh
var linPeas []byte

func LinEnum() *EnumScript { return &EnumScript{Data: linEnum} }
func LinPeas() *EnumScript { return &EnumScript{Data: linPeas} }
