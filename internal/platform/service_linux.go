//go:build linux

package platform

import "ProxyX/internal/platform/linux"

func newPlatformService() (Service, error) {
	return linux.New(), nil
}
