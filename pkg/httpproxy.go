package pkg

import (
	"net/http"
	"strconv"

	"github.com/gtt116/hass/log"
)

type HTTPProxy struct {
	IPAddr string
	Port   int
	Proxy
}

func (self *HTTPProxy) Serve() error {
	listenAddr := self.IPAddr + ":" + strconv.Itoa(self.Port)
	log.Infof("HTTP proxy listen at: %v", listenAddr)
	http.ListenAndServe(listenAddr, self)
	return nil
}

func (self *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("HTTP %v %v\n", r.Method, r.Host)

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var target *Target

	if r.Method == "CONNECT" {
		conn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
		target, err = NewTarget(r.Host, conn)
		if err != nil {
			log.Debugf("%v\n", err)
		}
	} else {
		r.Header.Del("Proxy-Connection")
		r.Header.Del("Proxy-Authenticate")
		r.Header.Del("Proxy-Authorization")
		// Connection, Authenticate and Authorization are single hop Header:
		// http://www.w3.org/Protocols/rfc2616/rfc2616.txt
		// 14.10 Connection
		//   The Connection general-header field allows the sender to specify
		//   options that are desired for that particular connection and MUST NOT
		//   be communicated by proxies over further connections.
		r.Header.Del("Connection")

		target, err = NewTarget(r.Host, conn)
		if err != nil {
			log.Debugf("%v\n", err)
		}

		// To proxy old http request, we need the raw request object.
		target.req = r
	}

	if err := self.Proxy.DoProxy(target); err != nil {
		log.Debugf("HTTP Proxy CONNECT failed: %v", err)
	}
}
