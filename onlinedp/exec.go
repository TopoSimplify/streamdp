package onlinedp

import "bytes"

func (self *OnlineDP) ExecuteTransaction(queries []string) {
		var buf bytes.Buffer
		buf.WriteString("BEGIN;\n")
		for _, q := range queries {
			buf.WriteString(q + "\n")
		}
		buf.WriteString("COMMIT;\n")
		var trans = buf.String()
		_, err := self.Src.Exec(trans)
		if err != nil {
			panic(err)
		}
}
