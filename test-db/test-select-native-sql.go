package main

import "log"

func TestSelectNativeSQL() {
	log.Println(">>> TestSelectNativeSQL() ----------------------------------------------------------------------")

	if err := db.Client1.Ping(); err != nil {
		log.Fatal(err)
	}

	var (
		id   int
		name string
	)

	rows, err := db.Client1.Query(selectUserQuery, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
