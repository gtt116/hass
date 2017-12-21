package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gtt116/hass/log"
	"github.com/nu7hatch/gouuid"
)

func probeUrl() string {
	uuid4, _ := uuid.NewV4()
	return fmt.Sprintf("http://speed.open.ad.jp:8080/download?nocache=%v&size=500000", uuid4.String())
}

func check(config *Config) {
	defer time.Sleep(time.Second * 10)

	proxyUrl, err := url.Parse(fmt.Sprintf("socks5://127.0.0.1:%d", config.Local.SocksPort+1))
	if err != nil {
		log.Errorln("probe error:", err)
		return
	}
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	resp, err := myClient.Get(probeUrl())
	if err != nil {
		log.Errorln("probe error:", err)
		return
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln("probe error:", err)
		return
	}
}

func StartChecker(config *Config) {
	for {
		check(config)
	}
}
