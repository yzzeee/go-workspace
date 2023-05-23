package main

/*
https://jmoiron.github.io/sqlx/
http://go-database-sql.org/importing.html

mysql> create user 'hello-user'@'%' identified by '1234';
Query OK, 0 rows affected (0.02 sec)

mysql> grant all on *.* to 'hello-user'@'%';
Query OK, 0 rows affected (0.01 sec)

mysql> flush privileges;
Query OK, 0 rows affected (0.01 sec)

CREATE TABLE `user` (
    `id` INT(10)  NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(20) NOT NULL
);

DROP TABLE `user`;

INSERT INTO `hello-mysql`.`user` (name) VALUES('홍길동');
INSERT INTO `hello-mysql`.`user` (name) VALUES('둘리');
INSERT INTO `hello-mysql`.`user` (name) VALUES('뽀로로');
INSERT INTO `hello-mysql`.`user` (name) VALUES('뽀로로');
*/

var db *DB

func init() {
	var err error
	db, err = NewDB("hello-mysql", "hello-user", "1234", "localhost:3306")
	if err != nil {
		panic(err)
	}
}

// 쿼리
var (
	selectUserQuery       = "SELECT id, name FROM `user` WHERE id = ?"
	upsertWidgetListQuery = `INSERT INTO
							   dashboard_widget (user_id, i, x, y, w, h, widget_code, side_information)
						     VALUES
							   (:user_id, :i, :x, :y, :w, :h, :widget_code, :side_information)
						     ON DUPLICATE KEY
						     UPDATE
							   x=:x, y=:y, w=:w, h=:h, widget_code=:widget_code, side_information=:side_information`
	insertEventQuery         = `INSERT INTO events(name, properties, browser) values (?, ?, ?)`
	selectEventByIdQuery     = `SELECT * FROM events WHERE id = ?`
	deleteWidgetListQuery    = `DELETE FROM dashboard_widget WHERE i NOT IN (?) and x = ?`
	deleteAllWidgetListQuery = `DELETE FROM dashboard_widget`
)

func main() {
	defer db.Client1.Close()
	defer db.Client2.Close()
	TestSelectNativeSQL()
	TestSelectSQLx()
	TestStringInterfaceMap()
	TestUpsert()
	TestDelete()
}
