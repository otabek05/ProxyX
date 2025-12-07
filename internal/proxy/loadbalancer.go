package proxy

import (
	"ProxyX/internal/common"
	"net/url"
	"sync"
)


type LoadBalancer struct {
	backends []*url.URL
	index int
	mutex sync.Mutex
}

func NewLoadBalancer(backendUrls []common.ProxyServer) (*LoadBalancer, error) {
	var backends []*url.URL

	for _, addr := range backendUrls {
		u, err := url.Parse(addr.URL)
		if err != nil {
			return nil, err 
		}

		backends = append(backends, u)
	}

	return &LoadBalancer{
		backends: backends,
	}, nil
}


func (l *LoadBalancer) Next() *url.URL {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if len(l.backends) == 0 {
		return nil
	}

	target := l.backends[l.index%len(l.backends)]
	l.index++
	return target
}


