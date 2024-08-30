package banlancer

import "errors"

var (
	NoHostError              = errors.New("no host")
	AlgorithmNotSupportError = errors.New("algorithm not supported")
)

// Balancer 处理反向代理的平衡器接口
type Balancer interface {
	Add(string)
	Remove(string)
	Balance(string) (string, error)
	Inc(string)
	Done(string)
}

// Factory 反向代理的平衡器工厂
type Factory func([]string) Balancer

var factories = make(map[string]Factory)

// Build 根据算法名称创建一个负载均衡器
func Build(algorithm string, host []string) (Balancer, error) {
	if f, ok := factories[algorithm]; ok {
		return f(host), nil
	}
	return nil, AlgorithmNotSupportError
}
