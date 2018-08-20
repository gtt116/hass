package pkg

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseConfigFile(t *testing.T) {

	file, err := ioutil.TempFile("", "hass")
	defer os.Remove(file.Name())

	data := `
local:
  socks_port: 7072
  http_port: 8888

backend:
  timeout: 5
  balance: "rr"  # rr, uri, least-conn
  port: 29123
  password: 123123
  method: "base64"
  iplist: www.gaott.info/hass.txt

  servers:
  - ip: 188.42.254.18
    port: 123213
    password: XXXX
    method: "rc4-md5"

  - ip: 27.0.232.109
`

	ioutil.WriteFile(file.Name(), []byte(data), 'w')

	c, err := ParseConfigFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if c.Local.SocksPort != 7072 {
		t.Error("socksPort should be 7072")
	}

	if c.Local.HttpPort != 8888 {
		t.Error("HttpPort should be 8888")
	}

	if c.Local.Host != "127.0.0.1" {
		t.Error("Host should be 127.0.0.1")
	}

	if c.Local.AdminPort != 7777 {
		t.Error("AdminPort should be 7777")
	}

	if c.Backend.Timeout != 5 {
		t.Error("Backend timeout should be 5")
	}

	if c.Backend.Balance != "rr" {
		t.Error("Backend balance should be rr")
	}

	if len(c.Backend.Servers) != 2 {
		t.Error("Should have 2 backend server")
	}

	if c.Backend.Servers[0].Method != "rc4-md5" {
		t.Error("Server 0 method should be rc4-md5")
	}

	if c.Backend.Servers[1].Port != 29123 {
		t.Error("Server 1 port should be 29123")
	}

	if c.Backend.Servers[1].Password != "123123" {
		t.Error("Server 1 password should be 123123")
	}

	if c.Backend.Servers[1].Method != "base64" {
		t.Error("Server 1 method should be base64")
	}

	if c.Backend.IPList != "www.gaott.info/hass.txt" {
		t.Error("IPList should be www.gaott.info/hass.txt")
	}
}
