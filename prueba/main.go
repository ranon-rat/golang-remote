package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type api struct {
	command string
}

func main() {
	comandos := `{"command":"caf"}`
	var m api
	j := json.NewDecoder(strings.NewReader(comandos))
	err := j.Decode(&m)
	if err != nil {

		fmt.Println(err)
	}
	fmt.Println(m)

}
