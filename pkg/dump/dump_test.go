package dump

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDumpTransport(t *testing.T) {

	assert := require.New(t)

	buf := &bytes.Buffer{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	server.Client().Transport = WrapTransport(buf, server.Client().Transport)

	req := request("POST", server.URL, bytes.NewBufferString("{}"))

	res, err := server.Client().Do(req)
	assert.NoError(err)
	assert.Equal(200, res.StatusCode)

	assert.Contains(buf.String(), "{}\n")
}

func request(method, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)

	return req
}
