package learn

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetPid(t *testing.T) {
	pid := os.Getpid()
	assert.True(t, pid > 0)
}
