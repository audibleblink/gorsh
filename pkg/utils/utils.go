package utils

import (
	"io/ioutil"
	"net/http"
)

func GetBytes(fs http.FileSystem, filename string) (data []byte, err error) {
	file, err := fs.Open(filename)
	if err != nil {
		return data, err
	}
	defer file.Close()

	data, err = ioutil.ReadAll(file)
	if err != nil {
		return data, err
	}
	return
}
