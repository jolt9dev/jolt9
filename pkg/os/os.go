package os

import (
	"fmt"
	"runtime"
)

const (
	PLATFORM = runtime.GOOS
	ARCH     = runtime.GOARCH
)

var (
	ErrOsNotSupported = fmt.Errorf("os %s not supported", runtime.GOOS)
)

func init() {
}
