package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetIpAddress(t *testing.T) {
	request := &http.Request{
		Header: map[string][]string{
			"X-Forwarded-For": []string{"27.68.135.90"},
		},
	}
	s := GetIPAddresses(request)
	assert.Equal(t, "27.68.135.90", s)
}
