// +build zstd

package zip

import (
	"github.com/valyala/gozstd"
	"io/ioutil"
)

// Bytes compresses the file contents of the given filepath string
func Bytes(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return file, err
	}

	zipData := gozstd.CompressLevel(nil, file, 15)
	return zipData, err
}
