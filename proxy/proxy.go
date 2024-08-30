package proxy

import (
	"MyLoadBanlancer/banlancer"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	XRealIP       = http.CanonicalHeaderKey("X-Real-IP")
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	XProxy        = http.CanonicalHeaderKey("X-Proxy")
)

var (
	ReverseProxy = "Balancer-Reverse-Proxy"
)

type HTTPProxy struct {
	hostMap map[string]*httputil.ReverseProxy
	lb      banlancer.Balancer

	sync.RWMutex
	alive map[string]bool
}

// NewHTTPProxy 创建新方向代理携带url和平衡器算法
func NewHTTPProxy(targetHosts []string, algorithm string) (*HTTPProxy, error) {
	lb, err := banlancer.Build(algorithm, targetHosts)
	if err != nil {
		return nil, err
	}
	hosts := make([]string, 0)
	hostMap := make(map[string]*httputil.ReverseProxy)
	alive := make(map[string]bool)
	for _, targetHost := range targetHosts {
		ul, err := url.Parse(targetHost)
		if err != nil {
			return nil, err
		}
		proxy := httputil.NewSingleHostReverseProxy(ul)

		orginDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			orginDirector(req)
			req.Header.Set(XRealIP, GetIP(req))
			req.Header.Set(XProxy, ReverseProxy)
		}
		host := GetHost(ul)
		alive[host] = true
		hostMap[host] = proxy
		hosts = append(hosts, host)
	}
	return &HTTPProxy{
		hostMap: hostMap,
		lb:      lb,
		alive:   alive,
	}, nil
}

func (h *HTTPProxy) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("proxy causes panic: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.(error).Error()))
		}
	}()
}
