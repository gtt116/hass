package pkg

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gtt116/hass/log"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

var (
	// key is shawdowsocks server addr like "192.2.2.2:8080"
	backends = make(map[string]*Backend)

	// Errors
	ErrNoAvailableServer = errors.New("No available server")
	ErrRingNotInit       = errors.New("Hash Ring not init")
)

type Backend struct {
	addr   string // id
	cipher *ss.Cipher

	InBytes        int64
	OutBytes       int64
	ConnCountCur   int64
	ConnCountTotal int64
	ConnCountErr   int64
	lock           sync.Mutex
}

func (b *Backend) String() string {
	return b.addr
}

func (b *Backend) AddInBytes(bytes int64) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.InBytes += bytes
}

func (b *Backend) AddOutBytes(bytes int64) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.OutBytes += bytes
}

func (b *Backend) IncreseConnCount() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.ConnCountCur++
	b.ConnCountTotal++
}

func (b *Backend) DecreseConnCount() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.ConnCountCur--
}

func (b *Backend) AddErr() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.ConnCountErr++
}

func (b *Backend) Cipher() *ss.Cipher {
	return b.cipher.Copy()
}

func ConfigBackend(cfg *Config) error {
	for _, server := range cfg.Backend.Servers {
		err := addBackend(server.IP, server.Port, server.Method, server.Password)
		if err != nil {
			return err
		}
	}
	return nil
}

func addBackend(host string, port int, method string, password string) (err error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	cipher, err := ss.NewCipher(method, password)
	if err != nil {
		return err
	}

	backend := &Backend{addr: addr, cipher: cipher}
	backends[addr] = backend

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

// Connect to a ss server, return BackendConn object and error through channel.
func connBackend(target *Target, backend *Backend, ch chan *BackendConn) {
	log.Debugf("Connect ss: %v, request %v", backend, target)
	// dial timeout is 3 seconds.
	ssConn, err := ss.Dial(target.Addr(), backend.String(), backend.Cipher())
	ch <- &BackendConn{backend, ssConn, err}
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
					log.Debugf("closing bad connection: %v", backConn)
					backConn.ssConn.Close()
				}
			}
		}
	}()
	return chConn
}

// Get a connection to a backend server. The caller should close the connection.
func GetConnection(config *Config, target *Target) (conn net.Conn, backend *Backend, err error) {
	ch := make(chan *BackendConn)
	for _, backend := range backends {
		go connBackend(target, backend, ch)
	}

	chConn := choiceConnection(ch)

	select {
	case backConn := <-chConn:
		return backConn.ssConn, backConn.back, nil
	case <-time.After(time.Second * 10):
		return nil, nil, ErrNoAvailableServer
	}
}

func GetAllBackends() map[string]*Backend {
	return backends
}
