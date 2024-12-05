package os_test

import (
	"runtime"
	"testing"

	"github.com/jolt9dev/jolt9/pkg/os"
	"github.com/stretchr/testify/assert"
)

func TestOs(t *testing.T) {
	assert.Equal(t, os.PLATFORM, runtime.GOOS)
}
