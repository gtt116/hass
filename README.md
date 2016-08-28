# HASS (HA shadowsocks)

[haes:]

Like haproxy, hass is used to distrubute requests to multiple shadowsocks backends by some balance algorithms.

## Main features

* Speak HTTP and sock5 protocol
* Support multiple balance algorithms(url_hash, random, etc)
* Support web based management and statictis
* Support automatic weight backends by network latency. (TODO)
