package onlinedp

import (
	"fmt"
	"bytes"
)

func (self *OnlineDP) ExecuteTransaction(queries []string) {
	var buf bytes.Buffer
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(buf.String())
			panic("size err")
		}
	}()

	buf.WriteString("BEGIN;\n")
	for _, q := range queries {
		buf.WriteString(q + "\n")
	}
	buf.WriteString("COMMIT;\n")

	var trans = buf.String()
	if _, err := self.Src.Exec(trans); err != nil {
		panic(err)
	}
}
