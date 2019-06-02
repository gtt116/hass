package pkg

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gtt116/hass/log"
	"strings"
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

func check(config *Config, first bool) {

	defer func() {
		if !first {
			waitTime := rand.Intn(60)
			log.Infof("[probe] Sleep %d seconds for the next checking\n", waitTime)
			time.Sleep(time.Second * time.Duration(waitTime))
		}
	}()

	updateBackend(config)

	probe(config)

}

func updateBackend(config *Config) {
	updateServer := config.Backend.IPList

	if updateServer == "" {
		return
	}

	response, err := http.Get(updateServer)
	if err != nil {
		log.Errorf("Error when updating server: %+v\n", err)
		return
	}

	if response.StatusCode != 200 {
		log.Errorf("Updating server: %s return: %s\n", updateServer, response.Status)
		return
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Errorf("Error when updating server: %+v\n", err)
		return
	}

	ips := strings.Split(string(content), "\n")

	for _, ip := range ips {
		if ip == "" {
			continue
		}
		addOneBackend(ip, config.Backend.Port, config.Backend.Method, config.Backend.Password)
	}
}

func probe(config *Config) {
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
		check(config, true)
	}

	for {
		check(config, false)
	}

	log.Errorln("checker goroute abnormally exit..")
}
