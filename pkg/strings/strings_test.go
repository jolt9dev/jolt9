package strings_test

import (
	"testing"

	"github.com/jolt9dev/jolt9/pkg/strings"
	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	assert.Equal(t, strings.TEST, "TEST")
}
