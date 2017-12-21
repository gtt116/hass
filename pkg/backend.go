package pkg

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gtt116/hass/log"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

var (
	// key is shawdowsocks server addr like "192.2.2.2:8080"
	backends    = make(map[string]*Backend)
	backendList = make([]*Backend, 0)

	// Errors
	ErrNoAvailableServer = errors.New("No available server")
	ErrRingNotInit       = errors.New("Hash Ring not init")
	rr_index             = 0 // the index to RoundRobin balance
)

type Backend struct {
	addr   string // id, form (Host:Port)
	cipher *ss.Cipher
	ConnStats
}

type ConnStats struct {
	bytes    int64
	duration time.Duration
	bps      float64 // Byte per second
}

func NewBackend(addr string, cipher *ss.Cipher) *Backend {
	return &Backend{addr: addr, cipher: cipher}
}

// Do connect to ss by given a target addr(formed "Host:Port")
func (b *Backend) Conn(addr string) (net.Conn, error) {
	log.Debugf("Connect ss: %v, request %v", b, addr)

	// dial timeout is 3 seconds.
	return ss.Dial(addr, b.String(), b.Cipher())
}

func (b *Backend) String() string {
	return b.addr
}

func (b *Backend) Bps() string {
	return fmt.Sprintf("%vKB, %.2fKB/s", b.bytes/1024, b.bps/1024)
}

func (b *Backend) Cipher() *ss.Cipher {
	return b.cipher.Copy()
}

func (b *Backend) UpdateStats(stats *ConnStats) {
	duration := stats.duration
	bytes := stats.bytes
	// caculate bps
	var bps float64
	seconds := duration / time.Second
	if seconds > 0 {
		bps = float64(bytes) / float64(seconds)
	}

	// update stats
	b.bps = (b.bps*float64(b.bytes) + bps*float64(bytes)) / (float64(b.bytes) + float64(bytes))
	if b.bytes < bytes {
		b.bytes = bytes
	}
}

func ConfigBackend(cfg *Config) error {
	for _, server := range cfg.Backend.Servers {
		err := addOneBackend(server.IP, server.Port, server.Method, server.Password)
		if err != nil {
			return err
		}
	}
	return nil
}

func addOneBackend(host string, port int, method string, password string) (err error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	cipher, err := ss.NewCipher(method, password)
	if err != nil {
		return err
	}

	backend := NewBackend(addr, cipher)
	backends[addr] = backend

	backendList = append(backendList, backend)

	log.Debugln("Add backend:", addr)
	return nil
}

type BackendConn struct {
	back   *Backend
	ssConn net.Conn
	err    error
}

func (bc *BackendConn) String() string {
	return fmt.Sprintf("%v #%v", bc.back, bc.err)
}

func (bc *BackendConn) Close() {
	if !bc.Error() {
		bc.ssConn.Close()
	}
}

func (bc *BackendConn) Backend() *Backend {
	return bc.back
}

func (bc *BackendConn) Error() bool {
	return bc.err != nil
}

// Pick up the first connection return through a channel, close the other connections.
func choiceConnection(conns chan *BackendConn) chan *BackendConn {
	chConn := make(chan *BackendConn)
	go func() {
		fired := false

		for _ = range backends {
			backConn := <-conns
			if backConn.err == nil {
				if !fired {
					chConn <- backConn
					fired = true
				} else {
					log.Debugf("closing unused connection: %v", backConn)
					backConn.ssConn.Close()
				}
			}
		}
	}()

	return chConn
}

func connBackend(target *Target, backend *Backend, ch chan *BackendConn) {
	log.Debugf("Connect ss: %v, request %v", backend, target)
	// dial timeout is 3 seconds.
	ssConn, err := ss.Dial(target.Addr(), backend.String(), backend.Cipher())
	ch <- &BackendConn{backend, ssConn, err}
}

// Get a connection to a backend server. The caller should close the connection.
func BestConnection(config *Config, target *Target) (conn net.Conn, backend *Backend, err error) {
	if len(backendList) == 0 {
		return nil, nil, ErrNoAvailableServer
	}

	best := backendList[0]
	for _, b := range backendList {
		if b.bps > best.bps {
			best = b
		}
	}
	ch := make(chan *BackendConn)
	go connBackend(target, best, ch)

	// FIXME: If return error, retry another backend
	backConn := <-ch
	return backConn.ssConn, backConn.back, nil
}

// Using round robin algorithm to probe the backend server's bps.
func ProbeConnection(config *Config, target *Target) (conn net.Conn, backend *Backend, err error) {
	timeout := time.After(time.Second * 20)
	for {
		backend := backendList[rr_index]
		log.Infof("Probe server: %v #%v", backend, backend.Bps())
		rr_index += 1
		if rr_index > len(backendList)-1 {
			rr_index = 0
		}
		ch := make(chan *BackendConn)
		go connBackend(target, backend, ch)

		select {
		case backConn := <-ch:
			if backConn.Error() {
				continue
			} else {
				return backConn.ssConn, backConn.back, nil
			}
		case <-timeout:
			return nil, nil, ErrNoAvailableServer
		}
	}
}

func GetAllBackends() map[string]*Backend {
	return backends
}
