// +build linux darwin freebsd !windows

package fetch

import (
	"io"
	"net/http"
	"os"
)

func Get(uri string, path string) (int64, error) {
	size, err := _downloadFile(uri, path)
	if err != nil {
		return 0, err
	}
	return size, nil
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
