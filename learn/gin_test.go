package learn

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGin(t *testing.T) {
	srv := gin.Default()
	err := srv.Run()
	assert.Nil(t, err)
}
