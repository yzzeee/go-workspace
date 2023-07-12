package z

import (
	"fmt"
)

var Message string
var Global = map[string]string{
	"hello": "map",
}

func init() {
	fmt.Println("z")
}
