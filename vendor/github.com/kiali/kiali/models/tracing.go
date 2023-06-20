package models

import "time"

type JaegerInfo struct {
	Enabled              bool     `json:"enabled"`
	Integration          bool     `json:"integration"`
	URL                  string   `json:"url"`
	NamespaceSelector    bool     `json:"namespaceSelector"`
	WhiteListIstioSystem []string `json:"whiteListIstioSystem"`
}

type TracingQuery struct {
	Start       time.Time
	End         time.Time
	Tags        map[string]string
	MinDuration time.Duration
	Limit       int
	Cluster     string
}
