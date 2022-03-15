package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"os"
)

func main() {
	bs, err := os.ReadFile("test.json")
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
		return
	}
	var m interface{}
	_ = json.Unmarshal(bs, &m)
	bs, _ = json.Marshal(m)

	value := gjson.Get(string(bs), `defs.0.items.#(deployName=="dp1-gin-demo")`)
	fmt.Println(value)

	str, _ := sjson.Set(string(bs), `defs.#.items.#(deployName=="dp1-gin-demo").deployCommand`, []string{"ok"})
	value = gjson.Get(str, "defs.0.items.0")
	fmt.Println(value)

	str, _ = sjson.Delete(string(bs), `defs.0.items.0.deployNodePorts`)
	value = gjson.Get(str, "defs.0.items.0")
	fmt.Println(value)

	m = true
	bs, _ = json.Marshal(m)
	fmt.Println(string(bs))
}
