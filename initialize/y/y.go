package y

import (
	"fmt"
	"go-practice/z"
)

var HelloY = map[string]bool{
	"1": true,
}

func init() {
	fmt.Println("y", z.Global)
}
