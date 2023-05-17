package proxy

import (
	"Balancer/balancer"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	XRealIP       = http.CanonicalHeaderKey("X-Real-IP")       //用来传递客户端的真实IP地址
	XProxy        = http.CanonicalHeaderKey("X-Proxy")         //标识请求是否通过代理服务器发送。
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For") //记录请求经过的代理服务器IP，X-Forwarded-For: <client>, <proxy1>, <proxy2>
	ReverseProxy  = "Balancer-Reverse-Proxy"
)

type HTTPProxy struct {
	hostMap map[string]*httputil.ReverseProxy
	lb      balancer.Balancer
	alive   map[string]bool
	sync.RWMutex
}

func NewHTTPProxy(targetHosts []string, algorithms string) (*HTTPProxy, error) {
	hosts := make([]string, 0)
	hostMap := make(map[string]*httputil.ReverseProxy)
	alive := make(map[string]bool)
	for _, targetHost := range targetHosts {
		u, err := url.Parse(targetHost)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		proxy := httputil.NewSingleHostReverseProxy(u)
		Director := proxy.Director
		proxy.Director = func(req *http.Request) {
			Director(req)
			req.Header.Set(XProxy, ReverseProxy)
			req.Header.Set(XRealIP, GetIP(req))
		}
		host := GetHost(u)
		alive[host] = true
		hostMap[host] = proxy
		hosts = append(hosts, host)
	}
	lb, err := balancer.Build(algorithms, hosts)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &HTTPProxy{
		hostMap: hostMap,
		lb:      lb,
		alive:   alive,
	}, nil
}

func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("proxy cause panic %s", err)
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(err.(error).Error()))
		}
	}()
	host, err := h.lb.Balance(GetIP(r))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(fmt.Sprintf("balance error %s", err.Error())))
		return
	}
	h.lb.Inc(host)
	defer h.lb.Done(host)
	h.hostMap[host].ServeHTTP(w, r)
}
