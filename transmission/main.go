package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	type Person struct {
		Name string `json:"name,omitempty"`
		Age  int
	}
	t1 := Person{
		Name: "톰",
		Age:  10,
	}
	t2 := Person{
		Name: "제리",
		Age:  5,
	}
	t3 := Person{
		Name: "",
		Age:  5,
	}

	// test omit empty
	m0, _ := json.Marshal(t3)
	fmt.Println(string(m0))

	t4 := make([]Person, 2)
	t4[0] = t1
	t4[1] = t2

	var t5 []Person

	m1, err := json.Marshal(t4)
	fmt.Println(m1, err)
	json.Unmarshal(m1, &t5)
	fmt.Println(t5)
}
