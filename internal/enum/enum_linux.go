//build !windows !darwin

package enum

import (
	log "github.com/sirupsen/logrus"
)

func LinEnum() *enumScript {
	data, err := scripts.Find("linenum.sh")
	if err != nil {
		log.Error(err)
	}

	return &enumScript{Data: data}
}
