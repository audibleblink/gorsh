//build +windows

package enum

import (
	log "github.com/sirupsen/logrus"

	"github.com/audibleblink/gorsh/internal/utils"
)

func Sherlock() *EnumScript {
	bytes, err := utils.GetBytes(scripts, "sherlock.ps1")
	if err != nil {
		log.Error(err)
	}

	return &EnumScript{Data: bytes}
}

func Jaws() *EnumScript {
	bytes, err := utils.GetBytes(scripts, "jaws.ps1")
	if err != nil {
		log.Error(err)
	}

	return &EnumScript{Data: bytes}
}

func PowerUp() *EnumScript {
	bytes, err := utils.GetBytes(scripts, "powerup.ps1")
	if err != nil {
		log.Error(err)
	}

	return &EnumScript{Data: bytes}
}
