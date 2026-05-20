//go:build windows

package cli

import (
	"ProxyX/internal/proxy"
	"log"

	"golang.org/x/sys/windows/svc"
)

const windowsServiceName = "proxyx"

func runServer(srv *proxy.ProxyServer) {
	isService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("Failed to detect service context: %v", err)
	}

	if !isService {
		if err := srv.Start(); err != nil {
			log.Fatalf("ProxyX exited: %v", err)
		}
		return
	}

	if err := svc.Run(windowsServiceName, &proxyService{srv: srv}); err != nil {
		log.Fatalf("Service handler exited: %v", err)
	}
}

type proxyService struct {
	srv *proxy.ProxyServer
}

func (p *proxyService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	const accepted = svc.AcceptStop | svc.AcceptShutdown

	status <- svc.Status{State: svc.StartPending}

	errCh := make(chan error, 1)
	go func() {
		errCh <- p.srv.Start()
	}()

	status <- svc.Status{State: svc.Running, Accepts: accepted}

	for {
		select {
		case err := <-errCh:
			if err != nil {
				log.Printf("ProxyX server returned: %v", err)
				status <- svc.Status{State: svc.Stopped}
				return false, 1
			}
			status <- svc.Status{State: svc.Stopped}
			return false, 0
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				status <- svc.Status{State: svc.StopPending}
				// fasthttp has no graceful shutdown wired here; exit the
				// process so SCM marks the service stopped promptly.
				return false, 0
			}
		}
	}
}
