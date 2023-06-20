package manager

import (
	"echo-server/proxy"
	urlpath "echo-server/util/url"
	_ "encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	_ "github.com/kiali/kiali/models"
	"github.com/labstack/echo/v4"
	"github.com/thoas/go-funk"
	"log"
	"net/http"
	"time"
)

var kialiServerAddress = "http://172.18.255.200/kiali"

type ErrorMessage struct {
	Message string `json:"message"`
}

func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// ProxyKialiServer 사용자의 요청을 Kiali 로 프록시 처리한 후 응답값을 처리하여 반환
func ProxyKialiServer(c echo.Context) error {
	fmt.Println("Start ProxyKialiServer")

	fmt.Println("http get", c.Request().RequestURI)
	fmt.Println("http get", c.Request().RequestURI)

	for a, b := range *proxy.Proxies {
		fmt.Println(a, b)
	}

	requestURL := c.Request().RequestURI
	requestMethod := c.Request().Method

	find := funk.Find(*proxy.Proxies, func(proxy proxy.Proxy) bool {
		pattern := urlpath.New(proxy.Pattern)
		_, patternOK := pattern.Match(requestURL)
		methodOK := proxy.Method == requestMethod
		return patternOK && methodOK
	})

	if find == nil {
		return nil // nil 반환 시
	}

	client := req.C().
		SetUserAgent("my-custom-client").
		SetTimeout(5 * time.Second)

	var response interface{}
	var errMsg ErrorMessage
	var err error
	var clientResponse *req.Response
	clientRequest := client.R()

	clientResponse, err = clientRequest.
		SetSuccessResult(&response).
		SetErrorResult(&errMsg).
		EnableDump().
		Get(kialiServerAddress + requestURL)

	//if find != nil {
	//	switch find.(proxy.Proxy).Type.(type) {
	//	case models.IstioConfigDetails:
	//		clientResponse, err = clientRequest.
	//			SetSuccessResult(&response).
	//			SetErrorResult(&errMsg).
	//			EnableDump().
	//			Get(kialiServerAddress + requestURL)
	//	}
	//}

	if err != nil {
		log.Println("error:", err)
		log.Println(clientResponse.Dump())
		c.Error(err)
	}
	if clientResponse.IsErrorState() {
		fmt.Println(errMsg.Message)
		c.Error(errors.New(errMsg.Message))
	}

	if clientResponse.IsSuccessState() {
		return c.JSON(http.StatusOK, response)
	}

	return nil
}
