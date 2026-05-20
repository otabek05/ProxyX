//go:build windows

package windows

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const ServiceName = "proxyx"

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Start() error {
	m, svcHandle, err := openService()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	defer svcHandle.Close()

	if err := svcHandle.Start(); err != nil {
		return fmt.Errorf("failed to start ProxyX: %w", err)
	}
	return nil
}

func (s *Service) Stop() error {
	m, svcHandle, err := openService()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	defer svcHandle.Close()

	status, err := svcHandle.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("failed to stop ProxyX: %w", err)
	}

	deadline := time.Now().Add(30 * time.Second)
	for status.State != svc.Stopped {
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for ProxyX to stop")
		}
		time.Sleep(300 * time.Millisecond)
		status, err = svcHandle.Query()
		if err != nil {
			return fmt.Errorf("query service status: %w", err)
		}
	}
	return nil
}

func (s *Service) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}
	return s.Start()
}

func (s *Service) Status() error {
	m, svcHandle, err := openService()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	defer svcHandle.Close()

	status, err := svcHandle.Query()
	if err != nil {
		return fmt.Errorf("query service status: %w", err)
	}

	if status.State != svc.Running {
		fmt.Println("ProxyX is not running")
		return fmt.Errorf("ProxyX is not running")
	}

	fmt.Println("ProxyX is running (Windows service)")
	fmt.Printf("    PID     : %d\n", status.ProcessId)
	return nil
}

func openService() (*mgr.Mgr, *mgr.Service, error) {
	m, err := mgr.Connect()
	if err != nil {
		return nil, nil, fmt.Errorf("connect to service manager: %w", err)
	}
	svcHandle, err := m.OpenService(ServiceName)
	if err != nil {
		m.Disconnect()
		return nil, nil, fmt.Errorf("open service %q (is it installed?): %w", ServiceName, err)
	}
	return m, svcHandle, nil
}
