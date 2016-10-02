package main

import (
	"errors"
	"math/rand"
	"net"
	"strconv"
	"sync"

	"github.com/serialx/hashring"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

const (
	DefaultWeight = 100
)

// key is shawdowsocks server addr like "192.2.2.2:8080"
var backends = make(map[string]*Backend)

// For GetBackendByURI
var backendRing *hashring.HashRing

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

func InitBackend(cfg *Config) error {
	for _, server := range cfg.Backend.Servers {
		err := AddBackend(server.IP, server.Port, server.Method, server.Password)
		if err != nil {
			Debugf("Init backend see failed: %v", err)
		}
	}

	hosts := backendKeys()
	hostWeights := make(map[string]int, len(hosts))
	for _, host := range hosts {
		hostWeights[host] = DefaultWeight
	}
	backendRing = hashring.NewWithWeights(hostWeights)
	return nil
}

func AddBackend(host string, port int, method string, password string) (err error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	cipher, err := ss.NewCipher(method, password)
	if err != nil {
		return err
	}

	backend := &Backend{addr: addr, cipher: cipher}
	backends[addr] = backend

	Debugln("Add backend:", addr)
	return nil
}

func backendKeys() []string {
	backendLen := len(backends)
	keys := make([]string, backendLen)
	i := 0
	for k := range backends {
		keys[i] = k
		i++
	}
	return keys
}

// Retry until find an available connecton. The caller should close the
// connection.
func ConnBackend(config *Config, target *Target) (conn net.Conn, backend *Backend, err error) {
	for {
		addr, backend, err := ChoiceBackend(config, target)
		if err != nil {
			return nil, nil, err
		}

		ssConn, err := ss.Dial(target.Addr(), addr, backend.Cipher())
		if err != nil {
			backendRing = backendRing.RemoveNode(addr)
			backend.AddErr()
			Debugf("Removing backend %v for reason: %v", addr, err)
			continue
		}
		return ssConn, backend, nil
	}
}

// Choice the correct backend by algorithms specified in config file.
// weight down the server which Dial() failed.
func ChoiceBackend(config *Config, target *Target) (addr string, backend *Backend, err error) {
	const (
		Random    = "random"
		Url       = "url_hash"
		LeaseConn = "lease_conn"
	)
	switch config.Backend.Balance {
	case Random:
		return GetBackendRandom()
	case Url:
		return GetBackendByURI(target.Addr())
	default:
		return GetBackendByURI(target.Addr())
	}
}

func GetAllBackends() map[string]*Backend {
	return backends
}

/*
 Get backend by random algorithm
*/
func GetBackendRandom() (addr string, backend *Backend, err error) {
	keys := backendKeys()
	randint := rand.Intn(len(backends))
	key := keys[randint]
	return key, backends[key], nil
}

/*
 Get backend by consistent hash algorithm on target URL.
*/
func GetBackendByURI(url string) (addr string, backend *Backend, err error) {
	if backendRing == nil {
		return "", nil, errors.New("backendRing not init.")
	}

	key, ok := backendRing.GetNode(url)
	if !ok || key == "" {
		return "", nil, errors.New("No sensable server")
	}

	return key, backends[key], nil
}
