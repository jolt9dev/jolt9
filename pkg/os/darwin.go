//go:build darwin || ios || tvos || watchos

package os

const (
	FAMILY = "darwin"
)

func IsDarwin() bool {
	return true
}
