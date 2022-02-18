package find

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func finderGen(pattern string, files *[]string) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) (err2 error) {
		if err != nil {
			fmt.Println(err)
			return
		}

		if info.IsDir() {
			return
		}

		reg, err2 := regexp.Compile(pattern)
		if err2 != nil {
			return
		}

		if reg.MatchString(info.Name()) {
			*files = append(*files, path)
		}
		return
	}
}

func Find(root, pattern string) (files []string, err error) {
	searcher := finderGen(pattern, &files)
	err = filepath.Walk(root, searcher)
	return
}
