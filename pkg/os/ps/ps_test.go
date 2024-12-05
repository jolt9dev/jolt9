package ps_test

import (
	"testing"

	"github.com/jolt9dev/jolt9/pkg/os"
	"github.com/jolt9dev/jolt9/pkg/os/ps"
	"github.com/stretchr/testify/assert"
)

func TestUserAndGroupIds(t *testing.T) {

	if os.IsWindows() {
		t.Skip("Skipping test on windows")
	}

	assert.Greater(t, ps.Uid(), -1)
	assert.Greater(t, ps.Gid(), -1)

	assert.Greater(t, ps.Egid(), -1)
	assert.Greater(t, ps.Euid(), -1)
}

func TestCwd(t *testing.T) {
	cwd, err := ps.Cwd()
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	assert.NotEmpty(t, cwd)
}

func TestProcessId(t *testing.T) {
	pid := ps.Pid()
	assert.Greater(t, pid, 0)

	ppid := ps.Ppid()
	assert.Greater(t, ppid, 0)
}

func TestWrite(t *testing.T) {
	b, err := ps.WriteString("test.txt")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	assert.Greater(t, b, 0)
}
