package main

import (
	"flag"

	"github.com/gtt116/hass/log"
	"github.com/gtt116/hass/pkg"
	"github.com/gtt116/hass/rlimit"
)

func main() {
	configFile := flag.String("config", "hass.yaml", "Hass default config file (yaml)")
	verbose := flag.Bool("verbose", false, "Default logging level is ERROR, change to DEBUG.")
	flag.Parse()

	if *verbose {
		log.EnableDebug()
	}

	config, err := pkg.ParseConfigFile(*configFile)
	if err != nil {
		log.Fatalln("Parse config file failed: ", err)
	}
	config.Report()

	rlimit.Setrlimit()
	pkg.ConfigBackend(config)

	proxy := pkg.NewProxyer(config)

	admin := pkg.NewProxyAdmin(config, proxy)
	go admin.StartSampling()
	go admin.ServeHTTP()

	httpp := &pkg.HTTPProxy{
		IPAddr: config.Local.Host,
		Port:   config.Local.HttpPort,
		Proxy:  proxy,
	}

	go httpp.Serve()

	socks := &pkg.Socks5{
		Ipaddr: config.Local.Host,
		Port:   config.Local.SocksPort,
		Proxy:  proxy,
	}

	socks.Serve()
}
