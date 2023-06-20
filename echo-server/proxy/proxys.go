package proxy

import (
	"echo-server/handler"
	"github.com/kiali/kiali/models"
	"github.com/labstack/echo/v4"
)

type HandlerFunc func(c echo.Context) error

type Proxy struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc HandlerFunc
	Type        interface{}
}

type proxies []Proxy

var Proxies = &proxies{
	{
		"IstioConfigDetails",
		"GET",
		"/api/namespaces/:namespace/istio/:objectType/:object",
		handler.IstioConfigDetails,
		models.IstioConfigDetails{},
	},
	{
		"allIstioConfigs",
		"GET",
		"/api/istio/config",
		handler.Healthz,
		"faf",
	},
}
