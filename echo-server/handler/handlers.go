package handler

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func Healthz(c echo.Context) error {
	c.Set("Hello", "World")
	return c.NoContent(http.StatusOK)
}

func IstioConfigDetails(c echo.Context) error {
	var bytes = "{\"a\":1,\"b\":2,\"c\":3,\"d\":{\"e\":{\"f\":4}}}"
	var b interface{}
	_ = json.Unmarshal([]byte(bytes), &b)
	log.Printf("%+v\n", b)
	return c.JSON(http.StatusOK, b)
}
