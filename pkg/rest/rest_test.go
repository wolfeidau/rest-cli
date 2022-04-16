package rest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoRequest(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("TEST", r.Header.Get("TEST"))
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	}))

	client := New(http.DefaultClient)

	buf := &bytes.Buffer{}
	err := client.DoRequest(http.MethodGet, server.URL, map[string]string{"TEST": "TEST"}, http.NoBody, buf)
	assert.NoError(err)
	assert.Equal("{}\n", buf.String())
}
