package main

import (
	"encoding/json"
	"fmt"
	"go-practice/common"
)

func main() {
	// using gjson
	//var foo = "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{},\"value\":[1653896944.095,\"80\"]}, {\"metric\":{},\"value\":[1653896944.095,\"80\"]}]}}"
	//fmt.Println(gjson.Get(foo, "status"))
	//fmt.Println(gjson.Get(foo, "data.result.#.value"))

	// native
	//a := make(map[string]map[string][]map[string]interface{})
	//_ = json.Unmarshal([]byte(foo), &a)
	//for i, ele := range a["data"]["result"] {
	//	fmt.Println(i, ele["value"])
	//}

	// using gjson
	var bytes = "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{\"namespace\":\"openshift-monitoring\"},\"value\":[1654736096.03,\"1.261144727299423\"]},{\"metric\":{\"namespace\":\"openshift-kube-apiserver\"},\"value\":[1654736096.03,\"0.5197757513776472\"]},{\"metric\":{\"namespace\":\"openshift-etcd\"},\"value\":[1654736096.03,\"0.35481809519734603\"]},{\"metric\":{\"namespace\":\"openshift-marketplace\"},\"value\":[1654736096.03,\"0.23072161485265724\"]},{\"metric\":{\"namespace\":\"neis\"},\"value\":[1654736096.03,\"0.09351117828365008\"]}]}}"
	//fmt.Println(gjson.Get(bytes, "data.result.#.metric"))
	//fmt.Println(gjson.Get(bytes, "data.result.#.value"))

	// native
	b := make(map[string]interface{})
	_ = json.Unmarshal([]byte(bytes), &b)
	//fmt.Println(b["data"].(map[string]interface{})["result"].([]interface{}))
	for _, ele := range b["data"].(map[string]interface{})["result"].([]interface{}) {
		fmt.Println(common.Get(ele, "metric.namespace"), common.Get(ele, "value"))
	}

}
