// +build linux,amd64 windows,amd64

package zip

import (
	"github.com/valyala/gozstd"
	"io/ioutil"
)

func Bytes(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return file, err
	}

	zipData := gozstd.CompressLevel(nil, file, 15)
	return zipData, err
}
