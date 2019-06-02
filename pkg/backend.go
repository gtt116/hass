package pkg

import (
	"container/ring"
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"

	"github.com/gtt116/hass/log"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

const (
	ringSize = 30
)

var (
	// The original order of backends
	backendList = make([]*Backend, 0)

	// The backends list which is ordered by bps desc
	orderedBackendList = make([]*Backend, 0)

	// the index to RoundRobin balance
	rr_index = 0

	// Errors
	ErrNoAvailableServer = errors.New("No available server")
	ErrDuplicated        = errors.New("Duplicated register server")
)

type Backend struct {
	addr   string // id, form (Host:Port)
	cipher *ss.Cipher
	ConnStats
	BpsRing
}

// Statistic about a connection
type ConnStats struct {
	// total send/receive bytes
	bytes int64

	// total time used to transfer
	duration time.Duration

	// Byte per second
	bps float64
}

// Each backend remember a ring of last bps to caculate the average bps.
type BpsRing struct {
	head *ring.Ring
	now  *ring.Ring
}

func (b *Backend) MeanBps() float64 {
	var sum, count float64
	x := b.head
	for {
		value := x.Value
		if value != nil {
			sum += value.(float64)
			count++
		}
		x = x.Next()
		if x == b.head {
			break
		}
	}
	if count > 0 {
		return sum / count
	} else {
		return 0
	}
}

func NewBackend(addr string, cipher *ss.Cipher) *Backend {
	b := &Backend{addr: addr, cipher: cipher}
	b.head = ring.New(ringSize)
	b.now = b.head
	return b
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
	// caculate bps
	var newBps float64
	seconds := stats.duration / time.Second
	if seconds > 0 {
		newBps = float64(stats.bytes) / float64(seconds)
	} else {
		newBps = 0
	}

	// Append BpsRing
	b.now.Value = newBps
	b.now = b.now.Next()
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

	if haveBackend(addr) {
		return ErrDuplicated
	}

	cipher, err := ss.NewCipher(method, password)
	if err != nil {
		return err
	}

	backend := NewBackend(addr, cipher)

	backendList = append(backendList, backend)
	orderedBackendList = append(orderedBackendList, backend)

	log.Infoln("Add backend:", addr)
	return nil
}

func haveBackend(addr string) bool {
	for _, back := range backendList {
		if back.addr == addr {
			return true
		}
	}
	return false
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

func connBackend(target *Target, backend *Backend) (net.Conn, error) {
	log.Debugf("Connect ss: %v, request %v", backend, target)
	// dial timeout is 3 seconds.
	ssConn, err := ss.Dial(target.Addr(), backend.String(), backend.Cipher())
	return ssConn, err
}

// Get a connection to a backend server. The caller should close the connection.
func BestConnection(config *Config, target *Target) (net.Conn, *Backend, error) {
	if len(backendList) == 0 {
		return nil, nil, ErrNoAvailableServer
	}

	topList := orderedBackendList[:]
	for _, b := range topList {
		conn, err := connBackend(target, b)
		if err == nil {
			return conn, b, nil
		}
	}
	return nil, nil, ErrNoAvailableServer

}

// Using round robin algorithm to probe the backend server's bps.
func ProbeConnection(config *Config, target *Target) (conn net.Conn, backend *Backend, err error) {
	tried := 0
	for {
		backend := backendList[rr_index]
		rr_index += 1
		if rr_index > len(backendList)-1 {
			rr_index = 0
		}
		conn, err := connBackend(target, backend)
		tried += 1
		if err != nil {
			log.Infof("Ping-> %v failed: %v", backend, err)

			// duration = 0 results in 0 bps.
			recvStats := &ConnStats{}
			recvStats.duration = 0
			backend.UpdateStats(recvStats)

			if tried == len(backendList) {
				return nil, nil, ErrNoAvailableServer
			} else {
				continue
			}
		} else {
			return conn, backend, nil
		}
	}
}

type SortedBacked []*Backend

func (sb SortedBacked) Len() int {
	return len(sb)
}

func (sb SortedBacked) Swap(a, b int) {
	sb[a], sb[b] = sb[b], sb[a]
}

func (sb SortedBacked) Less(a, b int) bool {
	return sb[a].MeanBps() < sb[b].MeanBps()
}

// Sort by average byte per second, greater Bps mean better performance.
func DoSort() {
	sort.Sort(sort.Reverse(SortedBacked(orderedBackendList)))
	log.Infoln("Top:")
	for _, b := range orderedBackendList {
		log.Infof("  %v\t%v", b, b.MeanBps())
	}
}
