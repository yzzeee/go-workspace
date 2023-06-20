package main

import (
	"echo-server/manager"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", manager.HelloWorld)
	// Kiali API Route Root
	kiali := e.Group("/api/servicemesh/kiali")
	// 모든 요청을 하나의 Endpoint 에서 관리
	kiali.Any("*", manager.ProxyKialiServer)

	e.Logger.Fatal(e.Start(":1323"))
}
