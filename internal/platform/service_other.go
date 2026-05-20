//go:build !linux && !darwin && !windows

package platform

import "fmt"

func newPlatformService() (Service, error) {
	return nil, fmt.Errorf("unsupported OS")
}
