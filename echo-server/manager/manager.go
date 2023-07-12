package manager

import (
	"echo-server/proxy"
	urlpath "echo-server/util/url"
	_ "encoding/json"
	_ "errors"
	"fmt"
	"github.com/imroc/req/v3"
	_ "github.com/kiali/kiali/models"
	"github.com/labstack/echo/v4"
	"github.com/thoas/go-funk"
	"io"
	"log"
	"net/http"
	"strings"
)

var KialiServerAddress = "http://192.168.49.100/kiali"
var KialiApiGroupPrefix = "/api/console/servicemesh"

type Error struct {
	Message string `json:"message"`
}

func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func setCookies(c echo.Context) {
	var bodyParams []byte
	if c.Request().Body != nil {
		bodyParams, _ = io.ReadAll(c.Request().Body)
	}
	requestURL := KialiServerAddress + "/api/authenticate"
	res, _ := req.C().R().
		EnableDump().
		SetContentType(c.Request().Header.Get("Content-Type")).
		SetBody(bodyParams).
		Post(requestURL)

	for _, cookie := range res.Cookies() {
		fmt.Println(cookie)
		c.SetCookie(cookie)
	}
}

// ProxyKialiServer 사용자의 요청을 Kiali 로 프록시 처리한 후 응답값을 처리하여 반환
func ProxyKialiServer(c echo.Context) error {
	log.Println("Start ProxyKialiServer")
	defer func() {
		log.Println("End ProxyKialiServer")
	}()

	// 세션 획득
	setCookies(c)

	// proxy handler 를 정의한 경우 정의한 handler 에서 사용자의 요청을 처리하여 응답함
	if find := funk.Find(*proxy.Proxies, func(proxy proxy.Proxy) bool {
		pattern := urlpath.New(proxy.Pattern)
		_, patternOK := pattern.Match("/api" + strings.TrimPrefix(c.Request().URL.Path, KialiApiGroupPrefix))
		methodOK := proxy.Method == c.Request().Method
		return patternOK && methodOK
	}); find != nil {
		log.Println("Run Proxy HandlerFunc")
		// 핸들러에 컨텍스트 위임
		_ = find.(proxy.Proxy).HandlerFunc(c)
	} else {
		log.Println("정의된 proxy 패턴 없음 kiali 로 전달하여 응답값 그대로 전달 할 것")
		requestURL := KialiServerAddress + "/api" + strings.TrimPrefix(c.Request().RequestURI, KialiApiGroupPrefix)
		log.Println("RequestURL: ", requestURL)
		client := req.C().R()
		var res *req.Response
		var err error

		switch c.Request().Method {
		case "GET":
			res, err = req.C().R().
				EnableDump().
				//SetSuccessResult(&result).
				Get(requestURL)
		case "POST":
			var bodyParams []byte
			if c.Request().Body != nil {
				bodyParams, _ = io.ReadAll(c.Request().Body)
			}
			res, err = client.
				EnableDump().
				SetContentType(c.Request().Header.Get("Content-Type")).
				SetBody(bodyParams).
				Post(requestURL)
		}

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		if res.IsErrorState() {
			return c.JSON(res.StatusCode, res.Status)
		}
		if res.IsSuccessState() {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)
			return c.String(res.StatusCode, string(body))
		}
	}
	return nil
}
