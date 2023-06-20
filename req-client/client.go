package main

import (
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/kiali/kiali/handlers"
	"log"
	"time"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func main() {
	// 사용자가 요청한 url
	requestUrl := "http://172.18.255.200/kiali/api/config"

	fmt.Println(requestUrl)

	client := req.C().
		SetUserAgent("my-custom-client").
		SetTimeout(5 * time.Second)

	client.DevMode()

	var publicConfig handlers.PublicConfig
	var errMsg ErrorMessage
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetSuccessResult(&publicConfig).
		SetErrorResult(&errMsg).
		EnableDump().
		Get("http://172.18.255.200/kiali/api/config")

	if err != nil {
		log.Println("error:", err)
		log.Println("raw content:")
		log.Println(resp.Dump())
		return
	}

	if resp.IsErrorState() {
		fmt.Println(errMsg.Message)
		return
	}

	if resp.IsSuccessState() {
		fmt.Printf("%s (%s)\n", publicConfig.GatewayAPIEnabled, publicConfig.Deployment)
		return
	}

	log.Println("unknown status", resp.Status)
	log.Println("raw content:")
	log.Println(resp.Dump())
}
