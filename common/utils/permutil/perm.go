package permutil

import (
	"os"
	"runtime"
)

func IAmAdmin() bool {
	switch runtime.GOOS {
	case "windows":
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		if err != nil {
			return false
		}
		return true
	default:
		if os.Getuid() == 0 {
			return true
		}

		if os.Geteuid() == 0 {
			return true
		}
	}
	return false
}
