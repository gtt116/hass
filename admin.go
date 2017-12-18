package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strconv"
	"time"
)

func ByteSize(args ...interface{}) string {
	const (
		_          = iota // ignore first value by assigning to blank identifier
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
		PB
		EB
		ZB
		YB
	)
	if len(args) != 1 {
		return "-"
	}

	y, ok := args[0].(int64)
	if !ok {
		return "NaN"
	}
	i := float64(y)

	switch {
	case i >= YB:
		return fmt.Sprintf("%.2fYB", i/YB)
	case i >= ZB:
		return fmt.Sprintf("%.2fZB", i/ZB)
	case i >= EB:
		return fmt.Sprintf("%.2fEB", i/EB)
	case i >= PB:
		return fmt.Sprintf("%.2fPB", i/PB)
	case i >= TB:
		return fmt.Sprintf("%.2fTB", i/TB)
	case i >= GB:
		return fmt.Sprintf("%.2fGB", i/GB)
	case i >= MB:
		return fmt.Sprintf("%.2fMB", i/MB)
	case i >= KB:
		return fmt.Sprintf("%.2fKB", i/KB)
	default:
		return fmt.Sprintf("%.2fB", i)
	}
}

type ProxyAdmin struct {
	cfg      *Config
	Proxy    *Proxyer
	Backends map[string]*Backend

	sessionRateCur int
	sessionRateMax int
	// sessionTotal = proxy.connTotal
	// sessionCur = proxy.connCount()
}

func (adm *ProxyAdmin) httpHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.New("hass template")
	tpl = tpl.Funcs(template.FuncMap{"bytesize": ByteSize})

	tpl, err := tpl.Parse(HTML)
	if err != nil {
		Debugln("Parse template failed:", err)
	}

	// Get latest backends
	adm.Backends = GetAllBackends()

	err = tpl.Execute(w, adm)
	if err != nil {
		Debugln("Execute template failed:", err)
	}
}

func (adm *ProxyAdmin) ServeHTTP() {

	http.HandleFunc("/", adm.httpHandler)

	addr := net.JoinHostPort(adm.cfg.Local.Host, strconv.Itoa(adm.cfg.Local.AdminPort))

	Infoln("Admin listen at:", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		Fatalln("HTTP failed:", err)
	}
}

func (adm *ProxyAdmin) StartSampling() {
	for {

		// wait a second
		time.Sleep(time.Second)
	}
}

func (adm *ProxyAdmin) BackendCount() int {
	return len(adm.Backends)
}

func (adm *ProxyAdmin) InBytesTotal() int64 {
	var total int64
	for _, b := range adm.Backends {
		total += b.InBytes
	}
	return total
}

func (adm *ProxyAdmin) OutBytesTotal() int64 {
	var total int64
	for _, b := range adm.Backends {
		total += b.OutBytes
	}
	return total
}
