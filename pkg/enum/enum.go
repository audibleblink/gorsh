package enum

import (
	"encoding/base64"

	"golang.org/x/text/encoding/unicode"
)

type EnumScript struct {
	Data []byte
}

func (e EnumScript) String() string {
	return string(e.Data)
}

func (e EnumScript) Base64() string {
	return base64.StdEncoding.EncodeToString(e.Data)
}

func (e EnumScript) UTF16LEB64() (string, error) {
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
