//go:build windows

package platform

import "ProxyX/internal/platform/windows"

func newPlatformService() (Service, error) {
	return windows.New(), nil
}
