package client

import (
	"net/http"
	"time"

	"github.com/scch94/Grecharge-gateway/config"
)

var Client http.Client

func InitHtppClient() {
	tr := &http.Transport{
		MaxIdleConns:        config.Config.Client.MaxIdleConns,
		MaxConnsPerHost:     config.Config.Client.MaxConnsPerHost,
		MaxIdleConnsPerHost: config.Config.Client.MaxConnsPerHost,
		IdleConnTimeout:     time.Duration(config.Config.Client.IdleConnTimeoutSeconds) * time.Second,
		DisableCompression:  config.Config.Client.DisableCompression,
		DisableKeepAlives:   config.Config.Client.DisableKeepAlives,
	}
	Client = http.Client{
		Transport: tr,
		Timeout:   time.Duration(config.Config.Client.PetitionsTimeOut) * time.Second,
	}
}
