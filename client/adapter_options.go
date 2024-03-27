package client

import (
	"net"
	url2 "net/url"
	"os"
	"strconv"
)

// tcp://127.0.0.1:11001?clientId=xxx&secret=xxx&connectNum=3
type AdapterOptions interface {
	Pid() string
	ClientId() string
	Ip() net.IP
	Port() int
	Secret() string
	ConnectNum() int
}

type adapterOptions struct {
	pid        string
	clientId   string
	ip         net.IP
	port       int
	secret     string
	connectNum int
}

func NewAdapterOptions(pid string, clientId string, ip net.IP, port int, secret string, connectNum int) AdapterOptions {
	return &adapterOptions{
		pid:        pid,
		clientId:   clientId,
		ip:         ip,
		port:       port,
		secret:     secret,
		connectNum: connectNum,
	}
}
func NewAdapterOptionsByUri(url string) AdapterOptions {
	var uri, err = url2.Parse(url)
	if err != nil {
		panic(err)
	}
	var port, err2 = strconv.Atoi(uri.Port())
	if err2 != nil {
		panic(err)
	}
	var connectNum = 1
	if uri.Query().Has("connectNum") {
		var p, e = strconv.Atoi(uri.Query().Get("connectNum"))
		if e != nil {
			panic(e)
		}
		connectNum = p
	}
	return &adapterOptions{
		pid:        uri.Query().Get("pid"),
		clientId:   uri.Query().Get("clientId"),
		ip:         net.ParseIP(uri.Host),
		port:       port,
		secret:     uri.Query().Get("secret"),
		connectNum: connectNum,
	}
}

func (a *adapterOptions) Ip() net.IP {
	var v, b = os.LookupEnv("VORTEX_ADAPTER_GATEWAY_HOST")
	if b {
		return net.ParseIP(v)
	}
	return a.ip
}

func (a *adapterOptions) Port() int {
	var v, b = os.LookupEnv("VORTEX_ADAPTER_GATEWAY_PORT")
	if b {
		var p, err = strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		return p
	}
	return a.port
}

func (a *adapterOptions) Secret() string {
	var v, b = os.LookupEnv("VORTEX_PRODUCT_SECRET")
	if b {
		return v
	}
	return a.secret
}

func (a *adapterOptions) ConnectNum() int {
	var v, b = os.LookupEnv("VORTEX_ADAPTER_CONNECTION_NUM")
	if b {
		var p, err = strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		return p
	}
	return a.connectNum
}

func (a *adapterOptions) Pid() string {
	var v, b = os.LookupEnv("VORTEX_PRODUCT_ID")
	if b {
		return v
	}
	return a.pid
}

func (a *adapterOptions) ClientId() string {
	var v, b = os.LookupEnv("VORTEX_ADAPTER_CLIENT_ID")
	if b {
		return v
	}
	return a.clientId
}
