//build +windows

package enum

import (
	log "github.com/sirupsen/logrus"
)

func Sherlock() *enumScript {
	data, err := scripts.Find("sherlock.ps1")
	if err != nil {
		log.Error(err)
	}

	return &enumScript{Data: data}
}

func Jaws() *enumScript {
	data, err := scripts.Find("jaws.ps1")
	if err != nil {
		log.Error(err)
	}

	return &enumScript{Data: data}
}

func PowerUp() *enumScript {
	data, err := scripts.Find("powerup.ps1")
	if err != nil {
		log.Error(err)
	}

	return &enumScript{Data: data}
}
