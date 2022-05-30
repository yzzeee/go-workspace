package main

import "fmt"

func main() {
	var queries = []string{"1: %s", "2: %s"}
	fmt.Println(queries)

	requests := make([]interface{}, len(queries))
	for i, query := range queries {
		requests[i] = query
	}
}
