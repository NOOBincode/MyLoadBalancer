package banlancer

import "sync"

type BaseBalancer struct {
	sync.RWMutex
	hosts []string
}

// Add 向负载平衡器添加一个新域名
func (b *BaseBalancer) Add(host string) {
	b.Lock()
	defer b.Unlock()
	for _, h := range b.hosts {
		if h == host {
			return
		}
	}
	b.hosts = append(b.hosts, host)
}

func (b *BaseBalancer) Remove(host string) {
	b.Lock()
	defer b.Unlock()
	for i, h := range b.hosts {
		if h == host {
			b.hosts = append(b.hosts[:i], b.hosts[i+1:]...)
			return
		}
	}
}

func (b *BaseBalancer) Balance(key string) (string, error) {
	b.RLock()
	defer b.RUnlock()
	if len(b.hosts) == 0 {
		return "", NoHostError
	}
	return b.hosts[0], nil
}
func (b *BaseBalancer) Inc(host string) {
	b.Add(host)
	return
}
func (b *BaseBalancer) Done(host string) {
	b.Remove(host)
}
