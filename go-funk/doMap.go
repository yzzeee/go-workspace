package main

import (
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"reflect"
)

type Widget struct {
	UserID     string `json:"userID" db:"user_id"`
	I          string `json:"i" db:"i"`
	X          int    `json:"x" db:"x"`
	Y          int    `json:"y" db:"y"`
	W          int    `json:"w" db:"w"`
	H          int    `json:"h" db:"h"`
	WidgetCode string `json:"widgetCode" db:"widget_code"`
}

func doMap() {
	var bytes = "[{\"user_id\":\"hell-user\",\"i\":\"123\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1231\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1232\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1233\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}},{\"user_id\":\"hello-user\",\"i\":\"1234\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widget_code\":\"TEXT_LABEL_WIDGET\",\"side_information\":{\"text\":\"test\"}}]"
	b := []Widget{}

	_ = json.Unmarshal([]byte(bytes), &b)
	fmt.Printf("%+v\n", b)

	widgetIds := funk.Map(b, func(widget Widget) string {
		return widget.I
	}).([]string)

	fmt.Println(reflect.TypeOf(widgetIds))
}
