package proxy

import (
	"ProxyX/internal/common"
	"ProxyX/internal/healthchecker"
	"net/url"
	"sync"
)


type LoadBalancer struct {
	servers []*healthchecker.Server
	index int
	mutex sync.Mutex
}



func NewLoadBalancer(backendUrls []common.ProxyServer, proxyConfig *common.ProxyConfig) (*LoadBalancer, error) {
	var servers []*healthchecker.Server

	for _, addr := range backendUrls {
		customUrl, err := url.Parse(addr.URL)
		if err != nil {
			return nil,err
		}

	    server :=  &healthchecker.Server{
			URL: customUrl,
			Path: proxyConfig.HealthCheck.Path,
			Health: true,
		}

		if proxyConfig.HealthCheck.Enabled {
			 if err := healthchecker.RegisterServer(addr.URL, server); err != nil {
				return nil, err 
			 }
		}

		servers = append(servers, server)
	}

	return &LoadBalancer{
		servers: servers,
	}, nil
}


func (l *LoadBalancer) Next() *url.URL {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	n := len(l.servers)
	if n == 0 {
		return nil
	}

	for range n {
	    srv := l.servers[l.index%n]
		l.index++
		if srv.IsHealthy(){
			return srv.URL
		}
	}

	return nil
}

