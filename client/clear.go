package main


func clearHistory(address string) {
	if _, err := post(address, []byte(`{"message":0}`)); err != nil {
		panic(err)
	}
}

