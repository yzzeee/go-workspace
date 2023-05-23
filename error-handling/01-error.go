package main

import (
	"fmt"
)

var (
	innerError = &innerErr{Msg: "innerMsg"}
)

// Custom error
type innerErr struct {
	Msg string
}

func (i *innerErr) Error() string {
	return "innerError"
}

// Functions
func innerFunc() error {
	return innerError
}

func middleFunc() error {
	if err := innerFunc(); err != nil {
		return fmt.Errorf("middleError")
	}
	return nil
}

func outerFunc() error {
	if err := middleFunc(); err != nil {
		return fmt.Errorf("outerError")
	}
	return nil
}

func doError() {
	// Get a error
	outerErr := outerFunc()

	// Print a error
	fmt.Printf("error: %v\n", outerErr)
}
