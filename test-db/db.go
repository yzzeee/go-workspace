package main

import (
	"database/sql"
	"fmt"
	// 드라이버를 import 해야 사용됨
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Client1 *sql.DB
	Client2 *sqlx.DB
}

func NewDB(dbName, dbUser, dbPassword, conn string) (*DB, error) {
	datasourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&timeout=60s&readTimeout=60s&charset=utf8",
		dbUser,
		dbPassword,
		conn,
		dbName)
	var client1, _ = sql.Open("mysql", datasourceName)
	// Open 하면서 Ping 실행
	var client2, _ = sqlx.Connect("mysql", datasourceName)
	return &DB{Client1: client1, Client2: client2}, nil
}
