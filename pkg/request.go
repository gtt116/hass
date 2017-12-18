package pkg

import (
	"net"
	"net/http"
	"strconv"
	"strings"
)

// A target contains all the information of a request.
type Target struct {
	// maybe IP or FQDN
	Host string

	// request server's port, mostly is 80
	Port int

	// Connection from client
	Client net.Conn

	// Only for httpproxy, I don't know how to make response directly to Client
	req *http.Request
}

// Parsing a "host:port" into a Target object
func NewTarget(hostStr string, client net.Conn) (*Target, error) {
	t := &Target{}

	var (
		port int
		err  error
	)

	parts := strings.Split(hostStr, ":")

	if len(parts) > 1 {
		port, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
	} else {
		port = 80
	}

	t.Host = parts[0]
	t.Port = port
	t.Client = client

	return t, nil
}

func (t *Target) Addr() string {
	return net.JoinHostPort(t.Host, strconv.Itoa(t.Port))
}
