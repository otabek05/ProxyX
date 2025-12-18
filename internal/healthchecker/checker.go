package healthchecker

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)


type Registry struct {
	servers []*Server
	mu sync.RWMutex
}

var global = &Registry{}


func RegisterServer(rawURL string, backend *Server) error {
	global.mu.Lock()
	global.servers = append(global.servers, backend)
	global.mu.Unlock()

	return nil
}

func Start(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func ()  {
		client := &http.Client{Timeout: 500 *time.Millisecond}
		for range ticker.C {
			global.mu.RLock()
			backends := global.servers
			global.mu.RUnlock()

			var wg sync.WaitGroup
			for _, backend :=  range backends {
				wg.Add(1)
				go func (b *Server)  {
					defer wg.Done()
					checkHealth(b, client)
				}(backend)
			}

		}	
	}()
}


func checkHealth(b *Server, client *http.Client) {
	u := b.URL.ResolveReference(&url.URL{Path: b.Path})
	resp , err := client.Get(u.String())
	health := err == nil && resp.StatusCode < 500
	b.SetHealthy(health)
}