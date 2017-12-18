# HASS (HA shadowsocks)

[haes:]

Like haproxy, hass is used to distribute requests to multiple shadowsocks backends by intelligent algorithms.

## Main features

* Speak HTTP and sock5 protocol
* Support web based management and statictis
* url persistent and lease connection first distribution.

## Demo

The admin web page look like:

![admin web](https://raw.githubusercontent.com/gtt116/hass/master/doc/hass.png)


## Build

```
git clone http://github.com/gtt116/hass
make
```

These commands will result in 3 executable file: hass, hass_mac, hass.exe

## Develop

if `govendor` do not existed, install it by `go get -u github.com/kardianos/govendor`, then
```
govendor build
```
