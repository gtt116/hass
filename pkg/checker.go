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

// Some webs from alexa top100.
var probeUrls = []string{
	"https://stackoverflow.com/",
	"https://www.google.com/",
	"https://www.youtube.com/",
	"https://www.facebook.com/",
	"https://twitter.com/",
	"https://outlook.live.com/",
	"https://www.netflix.com/",
	"https://www.wikipedia.org/",
	"https://www.wordpress.com/",
	"https://www.blogger.com",
	"https://www.apple.com/",
	"https://www.adobe.com/",
	"https://www.tumblr.com/",
	"https://www.amazon.com/",
	"https://vimeo.com/",
	"https://www.flickr.com/",
	"https://bitly.com/",
	"https://www.microsoft.com/zh-cn/",
	"https://godaddy.com/",
	"https://www.buydomains.com/",
	"https://www.reddit.com/",
	"https://www.w3.org/",
	"https://www.nytimes.com/",
	"http://europa.eu/",
	"http://statcounter.com/",
	"http://www.weebly.com/",
	"https://soundcloud.com/",
	"https://github.com/",
	"https://www.firefox.com",
	"https://en.gravatar.com/",
	"https://bluehost.com/",
}

func probeUrl() string {
	idx := rand.Intn(len(probeUrls))
	return probeUrls[idx]
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
