package enum

import (
	"encoding/base64"

	"github.com/gobuffalo/packr"
	"golang.org/x/text/encoding/unicode"
)

var scripts packr.Box

func init() {
	scripts = packr.NewBox("../../scripts")
}

type enumScript struct {
	Data []byte
}

func (e enumScript) String() string {
	return string(e.Data)
}

func (e enumScript) Base64() string {
	return base64.StdEncoding.EncodeToString(e.Data)
}

func (e enumScript) UTF16LEB64() (string, error) {
	return ToUnicode(string(e.Data))
}

func ToUnicode(in string) (string, error) {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	encoder := utf16.NewEncoder()
	bytes, err := encoder.Bytes([]byte(in))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), err
}
