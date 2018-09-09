// +build windows

package fetch

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
