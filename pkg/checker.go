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

func doProbe(config *Config, first bool) {
	defer func() {
		if !first {
			waitTime := rand.Intn(60)
			log.Infoln("[probe] Sleep %d seconds for the next checking", waitTime)
			time.Sleep(time.Second * time.Duration(waitTime))
		}
	}()

	proxyUrl, err := url.Parse("http://127.0.0.1:7073")
	if err != nil {
		log.Errorln("[probe] parse url error:", err)
		return
	}
	log.Debugln("[probe] sending probe request..")
	myClient := &http.Client{
		Timeout: time.Duration(10 * time.Second),
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
		doProbe(config, true)
	}

	for {
		doProbe(config, false)
	}

	log.Errorln("checker goroute abnormally exit..")
}
