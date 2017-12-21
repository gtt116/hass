package pkg

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Local   *LocalConfig
	Backend *BackendConfig
}

type LocalConfig struct {
	AdminPort int `yaml:"admin_port"`
	SocksPort int `yaml:"socks_port"`
	HttpPort  int `yaml:"http_port"`
	Host      string
}

type BackendConfig struct {
	Timeout  int // FIXME: delete it
	Balance  string
	Port     int
	Method   string
	Password string
	Servers  []*Server
}

type Server struct {
	IP       string
	Port     int
	Password string
	Method   string
}

func ParseConfig(data []byte) (config *Config, err error) {
	c := Config{}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	// sanity check
	if c.Backend == nil {
		return nil, errors.New("missing 'backend' section")
	}

	if c.Local == nil {
		return nil, errors.New("missing 'local' section")
	}

	if c.Local.SocksPort == 0 {
		return nil, errors.New("socks_port should not be 0")
	}

	// setup default port and password
	if c.Local.Host == "" {
		c.Local.Host = "127.0.0.1"
	}

	if c.Local.AdminPort == 0 {
		c.Local.AdminPort = 7777
	}

	if c.Backend.Method == "" {
		c.Backend.Method = "rc4-md4"
	}

	for _, server := range c.Backend.Servers {
		if server.Port == 0 {
			server.Port = c.Backend.Port
		}
		if server.Password == "" {
			server.Password = c.Backend.Password
		}
		if server.Method == "" {
			server.Method = c.Backend.Method
		}
	}

	return &c, nil
}

func ParseConfigFile(path string) (config *Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConfig(data)
}
