package zip

import (
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
)

// Bytes compresses the file contents of the given filepath string
func ZipWriter(path string) (compressed []byte, err error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var wr io.Writer
	zwr, err := zstd.NewWriter(wr)

	compressed = zwr.EncodeAll(file, compressed)
	return
}
