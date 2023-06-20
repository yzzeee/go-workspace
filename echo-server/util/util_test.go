package util

import (
	"echo-server/proxy"
	urlpath "echo-server/util/url"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/kiali/kiali/models"
	"github.com/thoas/go-funk"
	"log"
	"net/url"
	"strings"
	"testing"
	"time"
)

var kialiServerAddress = "http://172.18.255.200/kiali"

var prefix = "/api/console/servicemesh"

type ErrorMessage struct {
	Message string `json:"message"`
}

// api rules
var (
	/* NOTE: kiali 서버가 외부로 노출되어있나? -> No.. How to know...?
	   SECloudit 콘솔에서 클러스터별 Kiali API Address 를 알아야한다.
	   Kiali 도 깔아야 함, 어차피 Iframe 으로 보여주기로 했기 때문에 당연한 부분.
	   istio gateway 를 사용하기 위해서는 외부로 노출되어 있어야하고
	   kiali를 해당 gateway로 노출시키든지 어떤 방법으로든 서비스가 SE 에서 접속 가능해야함(like 프로메테우스) */
	istioConfigDetail = urlpath.New("/api/namespaces/:namespace/istio/:objectType/:object")
)

// TestUrlPathMatch 사용자의 요청 url 을 Match
func TestUrlPathMatch(t *testing.T) {
	// 요청 URL
	requestURL := "/api/console/servicemesh/namespaces/hello-namespace/istio/asfda%2Faf/dakfjkaf"

	// Prefix 제거
	requestURL = "/api" + strings.TrimPrefix(requestURL, prefix)

	t.Log("[REQUEST URL]", requestURL)
	match, ok := istioConfigDetail.Match(requestURL)
	if !ok {
		t.Error("[ERROR] url path dose not match")
	}

	namespace, _ := url.QueryUnescape(match.Params["namespace"])
	objectType, _ := url.QueryUnescape(match.Params["objectType"])
	t.Log(namespace, objectType)
}

func TestTrimStringPrefix(t *testing.T) {
	var s1 = "/api/servicemesh/kiali/adfjklajfldas?aaa=1&bbb=2"
	s1 = strings.TrimPrefix(s1, prefix)
	fmt.Println(s1)

	var s2 = "/api/servicemesh/kiali/adfjklajfldas"
	s2 = s2[len(prefix):]
	fmt.Println(s2)
}

func TestFindProxy(t *testing.T) {
	// 요청 URL
	//http://172.18.255.200/kiali/api/namespaces/istio-system/istio/envoyfilters/tcp-stats-filter-1.17?validate=true&help=true
	requestURL := "/api/console/servicemesh/namespaces/istio-system/istio/envoyfilters/tcp-stats-filter-1.17?validate=true&help=true"
	requestMethod := "POST"

	// Prefix 제거
	requestURL = "/api" + strings.TrimPrefix(requestURL, prefix)

	t.Log("[REQUEST_URL]", requestURL)

	// 정의된 프록시 목록에서 url 패턴과 메소드가 같은 것 찾기
	find := funk.Find(*proxy.Proxies, func(proxy proxy.Proxy) bool {
		pattern := urlpath.New(proxy.Pattern)
		_, patternOK := pattern.Match(requestURL)
		methodOK := proxy.Method == requestMethod
		return patternOK && methodOK
	})

	client := req.C().
		SetUserAgent("my-custom-client").
		SetTimeout(5 * time.Second)

	//client.DevMode()

	var errMsg ErrorMessage
	var err error
	var clientResponse *req.Response
	clientRequest := client.R()

	if find != nil {
		switch find.(proxy.Proxy).Type.(type) {
		case models.IstioConfigDetails:
			var istioConfigDetails models.IstioConfigDetails
			clientRequest.SetSuccessResult(&istioConfigDetails)
			defer func() {
				t.Log(333)
				if clientResponse.IsSuccessState() {
					t.Log(istioConfigDetails)
				}
			}()
		}
		clientResponse, err = clientRequest.
			SetErrorResult(&errMsg).
			EnableDump().
			Get(kialiServerAddress + requestURL)
		t.Log(111)
		if err != nil {
			log.Println("error:", err)
			log.Println(clientResponse.Dump())
			return
		}

		if clientResponse.IsErrorState() {
			fmt.Println(errMsg.Message)
			return
		}
	}
}
