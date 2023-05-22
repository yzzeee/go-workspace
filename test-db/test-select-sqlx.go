package main

import (
	"fmt"
	"log"
)

func TestSelectSQLx() {
	log.Println(">>> TestSelectSQLx() ---------------------------------------------------------------------------")
	type User struct {
		Id   string
		Name string
	}

	user := new([]User)
	err := db.Client2.Select(user, selectUserQuery, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
}
