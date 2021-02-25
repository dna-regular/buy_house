package http

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpGet(t *testing.T) {
	resp := HttpGet("http://www.baidu.com")
	log.Println(len(resp))
	assert.NotEmpty(t, resp)
}
