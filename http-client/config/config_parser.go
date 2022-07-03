package config

import (
	"github.com/Terry-Mao/goconf"
	"go-practice/http-client/kubernetes"
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

	ClientConfig.KubeConfigPath, err = configs.String("kube_config_path")
	if err != nil {
		panic(err)
	}
	kubernetes.KubeConfigPath = ClientConfig.KubeConfigPath

	ClientConfig.KubeIgnoreTLSVerification, err = configs.Bool("kube_ignore_tls_verification")
	if err != nil {
		panic(err)
	}
	kubernetes.IgnoreTLSVerification = ClientConfig.KubeIgnoreTLSVerification

	err = kubernetes.InitConfig()
	if err != nil {
		panic(err)
	}
}

// Init 환경 설정을 초기화 한다.
func Init() {
	err = conf.Parse(configLocation)
	if err != nil {
		panic(err)
	}

	parseClientConfig()
}
