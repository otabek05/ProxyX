//go:build darwin

package platform

import "ProxyX/internal/platform/darwin"

func newPlatformService() (Service, error) {
	return darwin.New(), nil
}
