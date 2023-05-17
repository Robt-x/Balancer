package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ConnectionTimeout = 3 * time.Second

func GetHost(u *url.URL) string {
	if _, _, err := net.SplitHostPort(u.Host); err == nil {
		return u.Host
	}
	if u.Scheme == "http" {
		return fmt.Sprintf("%s:%s", u.Host, "80")
	}
	if u.Scheme == "https" {
		return fmt.Sprintf("%s:%s", u.Host, "443")
	}
	return u.Host
}
func GetIP(req *http.Request) string {
	clientIP, _, _ := net.SplitHostPort(req.RemoteAddr)
	if len(req.Header.Get(XForwardedFor)) != 0 {
		xff := req.Header.Get(XForwardedFor)
		s := strings.Index(xff, ", ")
		if s == -1 {
			s = len(req.Header.Get(XForwardedFor))
		}
		clientIP = xff[:s]
	} else if len(req.Header.Get(XRealIP)) != 0 {
		clientIP = req.Header.Get(XRealIP)
	}
	return clientIP
}

func isBackendAlive(host string) bool {
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return false
	}
	resolveAddr := fmt.Sprintf("%s:%d", addr.IP, addr.Port)
	conn, err := net.DialTimeout("tcp", resolveAddr, ConnectionTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
