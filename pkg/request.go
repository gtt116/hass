package pkg

import (
	"net"
	"net/http"
	"strconv"
	"strings"
)

// A target contains all the information of a request.
type Target struct {
	Host    string // maybe IP or FQDN
	Port    int
	request *http.Request // only used by HTTP proxy
}

func (t *Target) Addr() string {
	return net.JoinHostPort(t.Host, strconv.Itoa(t.Port))
}

// parse a "host:port" into a Target object
func NewTarget(hostStr string) (*Target, error) {
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

	return t, nil
}
