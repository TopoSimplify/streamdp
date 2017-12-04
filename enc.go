package main

import (
	"log"
	"fmt"
	"bytes"
	"encoding/gob"
	"encoding/base64"
)

// go binary encoder
func Serialize(v *Vessel) string {
	var buf bytes.Buffer
	var err = gob.NewEncoder(&buf).Encode(v)
	if err != nil {
		log.Fatalln(err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// go binary decoder
func Deserialize(str string) *Vessel {
	var v *Vessel
	var dat, err = base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatalln(`failed base64 Decode`, err)
	}
	var buf bytes.Buffer
	_, err = buf.Write(dat)
	if err != nil {
		log.Fatalln(`failed to write to buffer`)
	}
	err = gob.NewDecoder(&buf).Decode(&v)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return v
}

