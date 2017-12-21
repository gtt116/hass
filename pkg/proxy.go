package pkg

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/gtt116/hass/log"
)

var (
	backStats map[string]*ConnStats
)

func init() {
	backStats = make(map[string]*ConnStats, 10)
}

type ConnTrack struct {
	LocalLocalAddr   string
	LocalRemoteAddr  string
	RemoteLocalAddr  string
	RemoteRemoteAddr string
	Target           string
	Latency          int64
}

type Proxy interface {
	DoProxy(target *Target) error
}

// Proxyer implement interface Proxy
// TODO: think about better name for it.
type Proxyer struct {
	cfg   *Config
	Probe bool // if probe is true, will do rr balance and record the bps of backend
}

func NewProxyer(config *Config) *Proxyer {
	return &Proxyer{
		cfg:   config,
		Probe: false,
	}
}

// Hass's version of socks5 server:
// pick up a backend shadowsocks server then pipe source and server.
func (p *Proxyer) DoProxy(target *Target) error {
	conn := target.Client
	defer conn.Close()

	startAt := time.Now()

	method := BestConnection
	if p.Probe {
		method = ProbeConnection
	}
	ssConn, backend, err := method(p.cfg, target)
	if err != nil {
		log.Errorf("▶ %v failed: %s", target, err)
		return err
	}
	defer ssConn.Close()

	latency := time.Since(startAt)
	log.Infof("▶ %v ▶ %v [%v]", backend, target, latency)

	// Maybe I can think about a way to make a valid request by my self :P
	if target.req != nil {
		target.req.Write(ssConn)
	}

	sendStats := &ConnStats{}
	recvStats := &ConnStats{}

	var wait sync.WaitGroup
	wait.Add(2)
	// client -> hass -> ss
	go copyStream(ssConn, conn, sendStats, &wait)
	// client <- hass <- ss
	go copyStream(conn, ssConn, recvStats, &wait)
	wait.Wait()

	if p.Probe {
		recvStats.duration += latency
		backend.UpdateStats(recvStats)
		DoSort()
	}
	return nil
}

// Copy stream from src to dst until EOF or some errors happend.
func copyStream(dst net.Conn, src net.Conn, stats *ConnStats, wait *sync.WaitGroup) {
	defer wait.Done()
	defer dst.Close()
	defer src.Close()

	start := time.Now()
	written, _ := io.Copy(dst, src)

	stats.bytes = written
	stats.duration = time.Since(start)
}
