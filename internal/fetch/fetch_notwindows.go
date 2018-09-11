// +build !windows

package fetch

func Get(uri string, path string) (int64, error) {
	size, err := _downloadFile(uri, path)
	if err != nil {
		return 0, err
	}
	return size, err
}
