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
	var bytes = "[{\"user_id\":\"hell-user\",\"i\":\"123\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1231\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1232\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1233\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1234\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}}]"
	b := []Widget{}

	// Unmarshal
	_ = json.Unmarshal([]byte(bytes), &b)
	fmt.Printf("%+v\n", b)

	widgetIds := funk.Map(b, func(widget Widget) string {
		return widget.I
	})

	tx, err := db.Client2.Beginx()
	query, args, err := sqlx.In(deleteWidgetListQuery, widgetIds)
	rows, err := db.Client2.Query(query, args...)
	fmt.Println(rows, err)
	tx.Commit()
}
