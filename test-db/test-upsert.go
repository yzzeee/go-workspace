package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Json map[string]interface{}

// Value ...
func (m Json) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return driver.Value(j), nil
}

func (m *Json) Scan(src interface{}) error {
	var source []byte
	_m := make(map[string]interface{})

	switch src.(type) {
	case []uint8:
		source = src.([]uint8)
	case nil:
		return nil
	default:
		return errors.New("incompatible type for Json")
	}
	err := json.Unmarshal(source, &_m)
	if err != nil {
		return err
	}
	*m = _m
	return nil
}

type Widget struct {
	UserID          string `json:"userID" db:"user_id"`
	I               string `json:"i" db:"i"`
	X               int    `json:"x" db:"x"`
	Y               int    `json:"y" db:"y"`
	W               int    `json:"w" db:"w"`
	H               int    `json:"h" db:"h"`
	WidgetCode      string `json:"widgetCode" db:"widget_code"`
	SideInformation Json   `json:"SideInformation" db:"side_information"`
}

/*
CREATE TABLE `dashboard_widget` (
  `user_id` varchar(45) NOT NULL COMMENT '사용자 아이디',
  `i` char(3) NOT NULL COMMENT '위젯 아이디',
  `x` int NOT NULL COMMENT 'X 좌표',
  `y` int NOT NULL COMMENT 'Y 좌표',
  `w` tinyint NOT NULL COMMENT '너비',
  `h` tinyint NOT NULL COMMENT '높이',
  `widget_code` varchar(45) NOT NULL COMMENT '위젯 코드',
  `side_information` longtext COMMENT '부가 정보(위젯의 추가 정보)',
  PRIMARY KEY(i)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='대시보드 위젯 목록';
*/

func TestUpsert() {
	log.Println(">>> TestUpsert() -------------------------------------------------------------------------------")
	// 트랜잭션 생성
	tx, err := db.Client2.Beginx()

	// 쿼리 준비
	namedStmt, err := tx.PrepareNamed(upsertWidgetListQuery)

	// JSON String
	var bytes = "{\"userID\":\"hello-user\",\"i\":\"123\",\"x\":1,\"y\":1,\"w\":1,\"h\":1,\"widgetCode\":\"TEXT_LABEL_WIDGET\",\"sideInformation\":{\"text\":\"test\"}}"
	b := Widget{}

	// Unmarshal
	_ = json.Unmarshal([]byte(bytes), &b)
	fmt.Printf("%+v\n", b)

	ss, err := json.MarshalIndent(b, " ", " ")
	fmt.Println("marshal indent", string(ss), err)
	result, err := namedStmt.Exec(b)

	fmt.Println(result, err)
	tx.Commit()
}
