package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/thoas/go-funk"
	"log"
)

func TestDelete() {
	log.Println(">>> TestDelete() -------------------------------------------------------------------------------")
	// JSON String
	//var bytes = "[{\"user_id\":\"hell-user\",\"i\":\"123\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1231\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1232\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1233\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1234\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}}]"
	var bytes2 = "[]"
	b := []Widget{}

	// Unmarshal
	_ = json.Unmarshal([]byte(bytes2), &b)
	fmt.Printf("%+v\n", b)

	widgetIds := funk.Map(b, func(widget Widget) string {
		return widget.I
	}).([]string)

	tx, _ := db.Client2.Beginx()
	if len(widgetIds) == 0 {
		rows, err := db.Client2.Exec(deleteAllWidgetListQuery)
		fmt.Println(rows, err)
	} else {
		query, args, err := sqlx.In(deleteWidgetListQuery, widgetIds)
		rows, err := db.Client2.Exec(query, args...)
		fmt.Println(rows, err)
	}
	tx.Commit()
}
