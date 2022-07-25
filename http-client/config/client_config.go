package config

// TODO: Change configuration location
var configLocation = "./.conf"

type clientConfig struct {
	PrometheusRequestURL      string `goconf:"default:prometheus_request_url"`       // PrometheusRequestURL: Request URL for prometheus
	PrometheusToken           string `goconf:"default:prometheus_token"`             // PrometheusToken: Token for prometheus
	KubeConfigPath            string `goconf:"default:kube_config_path"`             // KubeConfigPath: Path of the kube config
	KubeIgnoreTLSVerification bool   `goconf:"default:kube_ignore_tls_verification"` // KubeIgnoreTLSVerification: Ignore TLS verification when using Kubernetes API
}

// ClientConfig : clientConfig config structure
var ClientConfig clientConfig
