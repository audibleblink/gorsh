// +build windows !linux !darwin !freebsd

package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Get(uri string, path string) (int64, error) {
	var (
		err   error
		bytes int64
	)

	if uri[0:2] == "//" {
		bytes, err = _copy(uri, path)
	} else {
		bytes, err = _downloadFile(uri, path)
	}

	if err != nil {
		return 0, err
	}

	return bytes, err
}

func _copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func _downloadFile(uri string, path string) (int64, error) {
	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(uri)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return 0, err
	}

	info, _ := out.Stat()
	return info.Size(), nil
}
