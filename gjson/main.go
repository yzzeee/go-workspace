package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson" // like gson, jackson(at java)
)

func main() {
	// using gjson
	var foo = "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{},\"value\":[1653896944.095,\"80\"]}, {\"metric\":{},\"value\":[1653896944.095,\"80\"]}]}}"
	fmt.Println(gjson.Get(foo, "status"))
	fmt.Println(gjson.Get(foo, "data.result.#.value"))

	a := make(map[string]map[string][]map[string]interface{})
	_ = json.Unmarshal([]byte(foo), &a)
	for i, ele := range a["data"]["result"] {
		fmt.Println(i, ele["value"])
	}
}
