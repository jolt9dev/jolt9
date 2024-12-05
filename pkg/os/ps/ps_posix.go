//go:build aix || darwin || dragonfly || freebsd || hurd || illumos || ios || linux || netbsd || openbsd || plan9 || solaris || zos

package ps

import "os"

func IsElevated() bool {
	return os.Geteuid() == 0
}
