package proxy

import (
	"log"
	"time"
)

// ReadAlive 检查健康情况
func (h *HTTPProxy) ReadAlive(url string) bool {
	h.RLock()
	defer h.RUnlock()
	return h.alive[url]
}

// SetAlive 设置健康状态
func (h *HTTPProxy) SetAlive(url string, alive bool) {
	h.Lock()
	defer h.Unlock()
	h.alive[url] = alive
}

// HealthCheck 健康检查创建一个为每一个代理的goroutine 健康检测
func (h *HTTPProxy) HealthCheck(interval int) {
	for host := range h.hostMap {
		go h.healthcheck(host, interval)
	}
}

func (h *HTTPProxy) healthcheck(host string, interval int) {
	ticker := time.NewTimer(time.Duration(interval) * time.Second)
	for range ticker.C {
		if !IsBackendAlive(host) && h.ReadAlive(host) {
			log.Printf("%s is down", host)

			h.SetAlive(host, false)
			h.lb.Remove(host)
		} else if IsBackendAlive(host) && !h.ReadAlive(host) {
			log.Printf("site reachable, add %s to load balancer", host)
			h.SetAlive(host, true)
			h.lb.Add(host)
		}
	}
}
