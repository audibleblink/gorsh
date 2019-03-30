//build +windows

package enum

import (
	log "github.com/sirupsen/logrus"
)

func Sherlock() *EnumScript {
	data, err := scripts.Find("sherlock.ps1")
	if err != nil {
		log.Error(err)
	}

	return &EnumScript{Data: data}
}

func Jaws() *EnumScript {
	data, err := scripts.Find("jaws.ps1")
	if err != nil {
		log.Error(err)
	}

	return &EnumScript{Data: data}
}

func PowerUp() *EnumScript {
	data, err := scripts.Find("powerup.ps1")
	if err != nil {
		log.Error(err)
	}

	return &EnumScript{Data: data}
}
