package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gtt116/hass/log"
)

func allowOrigin(resp http.ResponseWriter) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
}

func hello(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hass is alive")
}

func total(resp http.ResponseWriter, req *http.Request) {
	allowOrigin(resp)

	type Total struct {
		Count       int64 `json:"count"`
		Connections int64 `json:"connections"`
		Sent        int64 `json:"sent"`
		Recv        int64 `json:"recv"`
	}

	t := &Total{0, 0, 0, 0}
	json.NewEncoder(resp).Encode(t)
}

func servers(resp http.ResponseWriter, req *http.Request) {
	allowOrigin(resp)

	type Server struct {
		Ip          string `json:"ip"`
		Sent        int64  `json:"sent"`
		Recv        int64  `json:"recv"`
		Connections int64  `json:"connections"`
		Msg         string `json:"msg"`
	}

	var servers []*Server
	servers = append(servers, &Server{Ip: "1.1.1.1"})
	servers = append(servers, &Server{Ip: "2.1.1.1"})
	servers = append(servers, &Server{Ip: "3.1.1.1"})
	servers = append(servers, &Server{Ip: "4.1.1.1"})

	json.NewEncoder(resp).Encode(servers)
}

func connections(resp http.ResponseWriter, req *http.Request) {
	allowOrigin(resp)

	type Connection struct {
		Source string `json:"source"`
		Server string `json:"server"`
		Target string `json:"target"`
		Sent   int64  `json:"sent"`
		Recv   int64  `json:"recv"`
	}
	var conns []*Connection

	conns = append(conns, &Connection{Source: "www.g.cn"})
	conns = append(conns, &Connection{Source: "www.g.cn"})
	conns = append(conns, &Connection{Source: "www.g.cn"})
	conns = append(conns, &Connection{Source: "www.g.cn"})

	json.NewEncoder(resp).Encode(conns)
}

func StartWebServer(config *Config) {
	http.Handle("/", http.FileServer(http.Dir("web/dist")))
	http.HandleFunc("/api/hello", hello)
	http.HandleFunc("/api/total", total)
	http.HandleFunc("/api/servers", servers)
	http.HandleFunc("/api/connections", connections)

	addr := fmt.Sprintf("%v:%v", config.Local.Host, config.Local.AdminPort)
	log.Infoln("Admin web listen at:", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
