package config

import (
	"github.com/Terry-Mao/goconf"
)

var conf = goconf.New()
var configs *goconf.Section
var err error

func parseClientConfig() {
	configs = conf.Get("default")
	if configs == nil {
		panic("failed to read config file")
	}

	ClientConfig = clientConfig{}
	ClientConfig.PrometheusRequestURL, err = configs.String("prometheus_request_url")
	if err != nil {
		panic(err)
	}

	ClientConfig.PrometheusToken, err = configs.String("prometheus_token")
	if err != nil {
		panic(err)
	}
}

func Init() {
	err = conf.Parse(configLocation)
	if err != nil {
		panic(err)
	}

	parseClientConfig()
}
