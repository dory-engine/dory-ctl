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

	value := gjson.Get(string(bs), `defs.0.items.0.builds`)
	fmt.Println(value)

	str, _ := sjson.Set(string(bs), `defs.0.items.0.builds.#(name=="dp1-gin-demo").run`, false)
	value = gjson.Get(str, `defs.0.items.0.builds`)
	fmt.Println(value)

	//str, _ = sjson.Delete(string(bs), `defs.0.items.0.deployNodePorts`)
	//value = gjson.Get(str, "defs.0.items.0")
	//fmt.Println(value)
	//
	//m = true
	//bs, _ = json.Marshal(m)
	//fmt.Println(string(bs))
}
