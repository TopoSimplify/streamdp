package onlinedp

import "encoding/base64"

func encode64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func decode64(b64str string) string {
	b, err := base64.StdEncoding.DecodeString(b64str)
	if err != nil {
		panic(err)
	}
	return string(b)
}
