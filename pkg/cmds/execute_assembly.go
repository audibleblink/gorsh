package cmds

import (
	"os"
)

var (
	hostname string
)

func GetHostname() string {
	if hostname != "" {
		return hostname
	}
	hostname, _ = os.Hostname()
	return hostname

}
