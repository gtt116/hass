package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gtt116/hass/log"
)

const (
	socks5Version = uint8(5)
	socks5nMethod = uint8(1)
	socks5Methods = uint8(0)
	cmdConnect    = uint8(1)
	socks5Reserve = uint8(0)

	ipv4Address = uint8(1)
	fqdnAddress = uint8(3)
	ipv6Address = uint8(4)
)

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

type DoProxy func(target *Target, conn net.Conn) error

type Socks5 struct {
	Ipaddr  string
	Port    int
	Handler DoProxy
}

func (s *Socks5) log(msg string) {
	fmt.Println(msg)
}

func (s *Socks5) Serve() error {
	listenAddr := s.Ipaddr + ":" + strconv.Itoa(s.Port)
	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	log.Infof("Socks5 listen at: %v", listenAddr)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Debugln("Accept error", err)
		}
		go s.serveConnection(conn)
	}
}

func (s *Socks5) serveConnection(conn net.Conn) error {
	defer conn.Close()

	var err error
	bufConn := bufio.NewReader(conn)

	err = s.negotiationAuth(bufConn, conn)
	if err != nil {
		log.Debugln("Negotiation failed", err)
		return err
	}

	target, err := s.parseTarget(bufConn, conn)
	if err != nil {
		log.Debugln("Get request failed", err)
		return err
	}

	if err = s.sendReply(conn); err != nil {
		log.Debugln("Send reply failed", err)
		return err
	}

	if s.Handler != nil {
		err = s.Handler(target, conn)
	} else {
		err = s.handleRequest(target, conn)
	}

	if err != nil {
		log.Debugln("Handler request failed", err)
	}

	return nil
}

/*
 Negotiation authentication method with client.

 clict frame:
	+----+----------+----------+
	|VER | NMETHODS | METHODS  |
	+----+----------+----------+
	| 1  |    1     | 1 to 255 |
	+----+----------+----------+

 server response with:

	+----+--------+
	|VER | METHOD |
	+----+--------+
	| 1  |   1    |
	+----+--------+

 METHOD:

	o  X'00' NO AUTHENTICATION REQUIRED
	o  X'01' GSSAPI
	o  X'02' USERNAME/PASSWORD
	o  X'03' to X'7F' IANA ASSIGNED
	o  X'80' to X'FE' RESERVED FOR PRIVATE METHODS
	o  X'FF' NO ACCEPTABLE METHODS

*/
func (s *Socks5) negotiationAuth(bufConn io.Reader, conn net.Conn) error {
	// version
	version := []byte{0}

	if _, err := bufConn.Read(version); err != nil {
		return err
	}

	if version[0] != socks5Version {
		err := fmt.Errorf("Unsupported SOCKS version: %v", version)
		return err
	}

	// n method
	nMethod := []byte{0}
	if _, err := bufConn.Read(nMethod); err != nil {
		return err
	}

	// method count
	methodCount := int(nMethod[0])

	methods := make([]byte, methodCount)
	if _, err := io.ReadAtLeast(bufConn, methods, methodCount); err != nil {
		return err
	}

	// version: 5, method: 0 (no authentication)
	if _, err := conn.Write([]byte{0x5, 0}); err != nil {
		return err
	}

	return nil
}

/*
The request frame:
	+----+-----+-------+------+----------+----------+
	|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+

*/
func (s *Socks5) parseTarget(bufConn io.Reader, conn net.Conn) (target *Target, err error) {

	header := []byte{0, 0, 0, 0}
	if _, err := io.ReadAtLeast(bufConn, header, len(header)); err != nil {
		return nil, err
	}

	if header[0] != socks5Version {
		return nil, fmt.Errorf("Unsupported command version: %v", header[0])
	}

	// TODO: support other methods
	if header[1] != cmdConnect {
		return nil, fmt.Errorf("Unsupported command method: %v", header[1])
	}

	if header[2] != socks5Reserve {
		return nil, fmt.Errorf("Unsupported command reserve: %v", header[2])
	}

	target = &Target{}

	switch header[3] {
	case ipv4Address:
		addr := make([]byte, 4)
		if _, err := io.ReadAtLeast(bufConn, addr, len(addr)); err != nil {
			return nil, err
		}
		target.Host = net.IP(addr).String()

	case ipv6Address:
		addr := make([]byte, 16)
		if _, err := io.ReadAtLeast(bufConn, addr, len(addr)); err != nil {
			return nil, err
		}
		target.Host = net.IP(addr).String()

	case fqdnAddress:
		buf := []byte{0}
		if _, err := bufConn.Read(buf); err != nil {
			return nil, err
		}
		fqdnLen := int(buf[0])
		fqdn := make([]byte, fqdnLen)
		if _, err := io.ReadAtLeast(bufConn, fqdn, fqdnLen); err != nil {
			return nil, err
		}
		target.Host = string(fqdn)

	default:
		return nil, fmt.Errorf("Unsupported command addrtype: %v", header[3])
	}

	// read port from frame
	port := make([]byte, 2)
	if _, err := io.ReadAtLeast(bufConn, port, len(port)); err != nil {
		return nil, err
	}
	target.Port = (int(port[0])<<8 | int(port[1]))

	return target, nil
}

// This is how the original socks5 server will do proxying
func (s *Socks5) handleRequest(target *Target, conn net.Conn) error {

	targetAddr := target.Addr()

	startAt := time.Now()
	proxyConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		return err
	}
	log.Debugf("Proxy %v (%v)", targetAddr, time.Since(startAt))

	errChan := make(chan error, 2)

	go CopyIO(conn, proxyConn, errChan)
	go CopyIO(proxyConn, conn, errChan)

	for i := 0; i < 2; i++ {
		e := <-errChan
		if e != nil {
			// return from this function closes target (and conn).
			return e
		}
	}

	return nil
}

/*
 always reply success.
*/
func (s *Socks5) sendReply(conn net.Conn) error {
	const (
		RespSuccess = uint8(0)
	)

	localIP := net.ParseIP("0.0.0.0").To4()
	localPort := uint16(33)

	// Format the message
	msg := make([]byte, 6+len(localIP))
	msg[0] = socks5Version
	msg[1] = RespSuccess
	msg[2] = socks5Reserve
	msg[3] = ipv4Address
	copy(msg[4:], localIP)
	msg[4+len(localIP)] = byte(localPort >> 8)
	msg[4+len(localIP)+1] = byte(localPort & 0xff)

	// Send the message
	_, err := conn.Write(msg)
	return err
}

func CopyIO(dst io.Writer, src io.Reader, errCh chan error) {
	_, err := io.Copy(dst, src)
	if tcpConn, ok := dst.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
	if err != nil {
		log.Debugf("Error on copy %v, err: %v", err)
	}
	errCh <- err
}
