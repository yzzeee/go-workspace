package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

/*
CREATE TABLE `events` (
  id int auto_increment primary key,
  name varchar(255),
  properties text,
  browser text
);
*/

type (
	StringInterfaceMap map[string]interface{}
	Event              struct {
		Id         int                `json:"id"`
		Name       string             `json:"name"`
		Properties StringInterfaceMap `json:"properties"`
		Browser    StringInterfaceMap `json:"browser"`
	}
)

func (m StringInterfaceMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return driver.Value(j), nil
}

func (m *StringInterfaceMap) Scan(src interface{}) error {
	var source []byte
	_m := make(map[string]interface{})

	switch src.(type) {
	case []uint8:
		source = src.([]uint8)
	case nil:
		return nil
	default:
		return errors.New("incompatible type for StringInterfaceMap")
	}
	err := json.Unmarshal(source, &_m)
	if err != nil {
		return err
	}
	*m = _m
	return nil
}

func insertEvent(db *sql.DB, event Event) (int64, error) {
	res, err := db.Exec(insertEventQuery, event.Name, event.Properties, event.Browser)
	if err != nil {
		return 0, err
	}
	lid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lid, nil
}

func selectEventById(db *sql.DB, id int64, event *Event) error {
	row := db.QueryRow(selectEventByIdQuery, id)
	err := row.Scan(&event.Id, &event.Name, &event.Properties, &event.Browser)
	if err != nil {
		return err
	}
	return nil
}

func buildPropertiesData() map[string]interface{} {
	return map[string]interface{}{
		"page": "/",
	}
}

func buildBrowserData() map[string]interface{} {
	return map[string]interface{}{
		"name": "Safari",
		"os":   "Mac",
		"resolution": struct {
			X int `json:"x"`
			Y int `json:"y"`
		}{1920, 1080},
	}
}

func TestStringInterfaceMap() {
	log.Println(">>> TestStringInterfaceMap() ----------------------------------------------------------------------")
	event := Event{
		Name:       "pageview",
		Properties: buildPropertiesData(),
		Browser:    buildBrowserData(),
	}

	insertedId, err := insertEvent(db.Client1, event)
	if err != nil {
		panic(err)
	}

	firstEvent := Event{}
	err = selectEventById(db.Client1, insertedId, &firstEvent)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nEvent fields:\n")
	fmt.Println("Id:         ", firstEvent.Id)
	fmt.Println("Name:       ", firstEvent.Name)
	fmt.Println("Properties: ", firstEvent.Properties)
	fmt.Println("Browser:    ", firstEvent.Browser)

	fmt.Println("\nJSON representation:\n")

	j, err := json.Marshal(firstEvent)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
