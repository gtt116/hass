package pkg

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gtt116/hass/log"
)

func probeUrl() string {
	return fmt.Sprintf("https://stackoverflow.com/")
}

func check(config *Config, first bool) {
	defer func() {
		if !first {
			time.Sleep(time.Second * time.Duration(rand.Intn(60)))
		}
	}()

	proxyUrl, err := url.Parse("http://127.0.0.1:7073")
	if err != nil {
		log.Errorln("[probe] parse url error:", err)
		return
	}
	myClient := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(proxyUrl),
			DisableKeepAlives: true,
		},
	}
	resp, err := myClient.Get(probeUrl())
	if err != nil {
		log.Errorln("[probe] send request error:", err)
		return
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("[probe] read response:", err)
		return
	}
}

func StartChecker(config *Config) {
	for _ = range backendList {
		check(config, true)
	}
	for {
		check(config, false)
	}
}
