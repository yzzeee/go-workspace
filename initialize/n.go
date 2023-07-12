package main

import (
	"fmt"
	"go-practice/initialize/k"
	"go-practice/initialize/y"
	"go-practice/z"
)

func init() {
	fmt.Println("n ->", z.Message, z.Global)
	fmt.Println("n", y.HelloY, k.A)
}
