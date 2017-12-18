package main

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/gtt116/hass/log"
)

func TimeAfter(timeout int) time.Time {
	return time.Now().Add(time.Second * time.Duration(timeout))
}

func CopyNetIO(dst net.Conn, src net.Conn, byteCh chan int64, msg string, timeout int) {
	written, err := CopyNetWithTimeout(dst, src, timeout)
	if tcpConn, ok := dst.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
	if err != nil {
		if nerr, ok := err.(net.Error); ok {
			if !nerr.Timeout() && !strings.Contains(nerr.Error(), "use of closed network connection") {
				// ignore i/o timeout error, and errClosing
				log.Errorf("net.Error %v : %v", msg, nerr)
			}
		} else {
			log.Errorf("Error %v: %v", msg, err)
		}
	}
	byteCh <- written
}

// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
func CopyNetWithTimeout(dst net.Conn, src net.Conn, timeout int) (written int64, err error) {
	buf := make([]byte, 32*1024)

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			src.SetReadDeadline(TimeAfter(timeout))
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				dst.SetWriteDeadline(TimeAfter(timeout))
			}
			if ew != nil { // Write failed
				err = ew
				dst.Close()
				src.Close()
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				dst.Close()
				src.Close()
				break
			}
		}
		if er == io.EOF {
			src.Close()
			break
		}
		if er != nil {
			src.Close()
			err = er
			break
		}
	}
	return written, err
}
