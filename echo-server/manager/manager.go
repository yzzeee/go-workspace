package manager

import (
	"echo-server/proxy"
	urlpath "echo-server/util/url"
	_ "encoding/json"
	_ "errors"
	"github.com/imroc/req/v3"
	_ "github.com/kiali/kiali/models"
	"github.com/labstack/echo/v4"
	"github.com/thoas/go-funk"
	"log"
	"net/http"
	"strings"
)

var kialiServerAddress = "http://172.18.255.200/kiali"
var kialiApiGroupPrefix = "/api/console/servicemesh"

type Error struct {
	Message string `json:"message"`
}

func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// ProxyKialiServer 사용자의 요청을 Kiali 로 프록시 처리한 후 응답값을 처리하여 반환
func ProxyKialiServer(c echo.Context) error {
	log.Println("Start ProxyKialiServer")
	defer func() {
		log.Println("End ProxyKialiServer")
	}()

	// proxy handler 를 정의한 경우 정의한 handler 에서 사용자의 요청을 처리하여 응답함
	if find := funk.Find(*proxy.Proxies, func(proxy proxy.Proxy) bool {
		pattern := urlpath.New(proxy.Pattern)
		_, patternOK := pattern.Match("/api" + strings.TrimPrefix(c.Request().URL.Path, kialiApiGroupPrefix))
		methodOK := proxy.Method == c.Request().Method
		return patternOK && methodOK
	}); find != nil {
		log.Println("Run Proxy HandlerFunc")
		// 핸들러에 컨텍스트 위임
		_ = find.(proxy.Proxy).HandlerFunc(c)
	} else {
		log.Println("정의된 proxy 패턴 없음 kiali 로 전달하여 응답값 그대로 전달 할 것")
		requestURL := kialiServerAddress + "/api" + strings.TrimPrefix(c.Request().RequestURI, kialiApiGroupPrefix)
		log.Println("RequestURL: ", requestURL)
		var result interface{}
		res, err := req.C().R().
			SetSuccessResult(&result).
			Get(requestURL)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		if res.IsErrorState() {
			return c.JSON(res.StatusCode, res.Status)
		}
		if res.IsSuccessState() {
			return c.JSON(res.StatusCode, result)
		}
	}
	return nil
}
