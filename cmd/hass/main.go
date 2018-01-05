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

	log.Infoln("====================> Hass <===========================")

	rlimit.Setrlimit()
	err = pkg.ConfigBackend(config)
	if err != nil {
		log.Fatalln("Config backend server failed: ", err)
	}

	proxy := pkg.NewProxyer(config)
	proxyProbe := pkg.NewProxyer(config)
	proxyProbe.Probe = true

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
	go socks.Serve()

	socksProbe := &pkg.HTTPProxy{
		IPAddr: "127.0.0.1",
		Port:   config.Local.ProbeHttpPort(),
		Proxy:  proxyProbe,
	}
	go socksProbe.Serve()

	go pkg.StartChecker(config)
	select {}
}
