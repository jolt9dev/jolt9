package os

import "runtime"

const (
	PLATFORM = runtime.GOOS
	ARCH     = runtime.GOARCH
)

func init() {
}
