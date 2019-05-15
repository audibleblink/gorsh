//build !windows !darwin

package enum

import (
	"github.com/audibleblink/gorsh/internal/utils"
)

func LinEnum() *EnumScript {
	bytes, err := utils.GetBytes(scripts, "linenum.sh")
	if err != nil {
		panic(err)
	}

	return &EnumScript{Data: bytes}
}
