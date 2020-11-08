package learn

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createTestServer() (r *gin.Engine, srv *httptest.Server) {
	r = gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "Jay")
	})
	r.Any("/echo", func(c *gin.Context) {
		raw, err := c.GetRawData()
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		switch accept := c.GetHeader("Accept"); accept {
		case "application/json":
			c.JSON(200, gin.H{
				"headers": c.Request.Header,
				"method":  c.Request.Method,
				"body":    raw,
			})
		default:
			c.Data(200, "application/echo", raw)
		}
	})
	srv = httptest.NewServer(r)
	return
}

func TestGin(t *testing.T) {
	_, ts := createTestServer()
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/hello")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Jay", string(body))
}

func TestDefer(t *testing.T) {
	val := 1
	plus1 := func(val *int) { *val += 1 }
	assertx := func(expect int, val *int) { assert.Equal(t, expect, *val) }
	defer assertx(3, &val) // "first statement executed last"
	defer plus1(&val)
	defer assertx(2, &val)
	defer plus1(&val)
}
