package main

import "fmt"

func main() {
	poo := []interface{}{"foo", "bar"}

	fmt.Println(fmt.Sprintf("%s, %s", poo...))
}
