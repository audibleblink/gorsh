package directory

import (
	"code.cloudfoundry.org/bytefmt"
	"io/ioutil"
)

func List(argv []string) (string, error) {
	var path string

	if len(argv) < 2 {
		path = "./"
	} else {
		path = argv[1]
	}

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return "", err
	}

	details := ""

	for _, f := range files {
		perms := f.Mode().String()
		size := bytefmt.ByteSize(uint64(f.Size()))
		modTime := f.ModTime().String()[0:19]
		name := f.Name()
		details = details + perms + "\t" + modTime + "\t" + size + "\t" + name + "\n"
	}

	return details, err
}
