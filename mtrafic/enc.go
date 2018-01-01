package mtrafic

import (
	"bytes"
	"encoding/gob"
	"encoding/base64"
)

// go binary encoder
func Serialize(v interface{}) (string, error) {
	var buf bytes.Buffer
	var err = gob.NewEncoder(&buf).Encode(v)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), err
}

// go binary decoder
func Deserialize(str string, v interface{}) (error) {
	var dat, err = base64.StdEncoding.DecodeString(str)
	if err != nil {
		return  err
	}
	var buf bytes.Buffer

	_, err = buf.Write(dat)
	if err != nil {
		return  err
	}

	err = gob.NewDecoder(&buf).Decode(v)
	if err != nil {
		return  err
	}
	return  err
}
