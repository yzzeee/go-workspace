package proxy

import (
	"echo-server/handler"
	"github.com/kiali/kiali/models"
	"github.com/labstack/echo/v4"
)

type handlerFunc func(c *echo.Context) error

type Proxy struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc handlerFunc
	Type        interface{}
}

type proxies []Proxy

var Proxies = &proxies{
	{
		"allIstioConfigs",
		"GET",
		"/api/istio/config",
		handler.Healthz,
		"faf",
	},
	{
		"IstioConfigDetails",
		"POST",
		"/api/namespaces/:namespace/istio/:objectType/:object",
		handler.IstioConfigDetails,
		models.IstioConfigDetails{},
	},
}
