// +build 386 darwin,amd64

package zip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Bytes(path string) ([]byte, error) {
	var zipData bytes.Buffer
	file, err := ioutil.ReadFile(path)
	if err != nil {
		bytes := zipData.Bytes()
		return bytes, err
	}

	zipper := gzip.NewWriter(&zipData)
	zipper.Write(file)
	zipper.Close()
	bytes := zipData.Bytes()
	return bytes, err
}
