//go:build !windows

package cli

import "os"

const elevationMessage = "This command must be run with sudo"

func isElevated() bool {
	return os.Geteuid() == 0
}
