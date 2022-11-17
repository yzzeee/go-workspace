package main

import (
	"fmt"
	"regexp"
)

func main() {
	match, _ := regexp.MatchString("^/api|^/auth", "/auth/test")
	fmt.Println(match)

	match2, _ := regexp.MatchString("/api/(.*)/auth/providers", "/api/stdfaf/auth/providers")
	fmt.Println(match2)

	ignorePaths := []string{"/api/(.*)/auth/providers"}

	var ignore = false
	for _, path := range ignorePaths {
		ignore, _ = regexp.MatchString(path, "/api/stdfaf/auth/providers")
	}

	fmt.Println(ignore)

}
