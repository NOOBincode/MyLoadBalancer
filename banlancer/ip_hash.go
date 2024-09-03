package banlancer

import "hash/crc32"

func init() {
	factories[IPHashBalancer] = NewIPHash
}

type IPHash struct {
	BaseBalancer
}

func NewIPHash(hosts []string) Balancer {
	return &IPHash{
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
	}
}

func (r *IPHash) Balance(key string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	value := crc32.ChecksumIEEE([]byte(key)) % uint32(len(r.hosts))
	return r.hosts[value], nil
}
