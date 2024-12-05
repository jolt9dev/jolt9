//go:build linux || android || solaris || illumos || plan9

package os

const (
	FAMILY = "linux"
)

func IsDarwin() bool {
	return false
}
