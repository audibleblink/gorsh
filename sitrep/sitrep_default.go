// +build linux darwin freebsd !windows

package sitrep

import (
	"net"
)

func Get(uri string, path string) (int64, error) {
	size, err := _downloadFile(uri, path)
	if err != nil {
		return 0, err
	}
	return size, nil
}
