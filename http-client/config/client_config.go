package config

// TODO: Change configuration location
var configLocation = "./.conf"

type clientConfig struct {
	PrometheusRequestURL string `goconf:"default:prometheus_request_url"` // PrometheusRequestURL: Request URL for prometheus
	PrometheusToken      string `goconf:"default:prometheus_token"`       // PrometheusToken: Token for prometheus
}

// ClientConfig : clientConfig config structure
var ClientConfig clientConfig
