//go:build windows
// +build windows

package fetch

func Get(uri string, path string) (int64, error) {
	var (
		err   error
		bytes int64
	)

	if uri[0:2] == "//" {
		bytes, err = Copy(uri, path)
	} else {
		bytes, err = _downloadFile(uri, path)
	}

	if err != nil {
		return 0, err
	}

	return bytes, err
}
