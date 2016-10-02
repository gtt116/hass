package main

import (
	"net/http"
	"strconv"
)

type HTTPProxy struct {
	IPAddr  string
	Port    int
	Handler DoProxy
}

func (self *HTTPProxy) Serve() error {
	listenAddr := self.IPAddr + ":" + strconv.Itoa(self.Port)
	Debugf("HTTP proxy listen at: %v", listenAddr)
	http.ListenAndServe(listenAddr, self)
	return nil
}

func (self *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Debugf("HTTP %v %v\n", r.Method, r.Host)

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
		target, err = NewTarget(r.Host)
		if err != nil {
			Debugf("%v\n", err)
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

		target, err = NewTarget(r.Host)
		if err != nil {
			Debugf("%v\n", err)
		}
		target.request = r
	}

	if err := self.Handler(target, conn); err != nil {
		Debugf("HTTP Proxy CONNECT failed: %v", err)
	}
}
