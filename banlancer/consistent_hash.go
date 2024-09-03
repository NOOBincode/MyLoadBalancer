package banlancer

import "github.com/lafikl/consistent"

func init() {
	factories[consistentHashBalancer] = NewConsistent
}

type Consistent struct {
	BaseBalancer
	ch *consistent.Consistent
}

func NewConsistent(hosts []string) Balancer {
	b := &Consistent{
		BaseBalancer: BaseBalancer{},
		ch:           consistent.New(),
	}
	for _, host := range hosts {
		b.Add(host)
	}
	return b
}

func (c *Consistent) Add(host string) {
	c.ch.Add(host)
}

func (c *Consistent) Remove(host string) {
	c.ch.Remove(host)
}

func (c *Consistent) Balance(key string) (string, error) {
	if len(c.ch.Hosts()) == 0 {
		return "", NoHostError
	}
	return c.ch.Get(key)
}
