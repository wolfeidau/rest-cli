package dump

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type DumpTransport struct {
	output    io.Writer
	transport http.RoundTripper
}

func WrapTransport(output io.Writer, transport http.RoundTripper) *DumpTransport {
	return &DumpTransport{
		output:    output,
		transport: transport,
	}
}

func (ds *DumpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	data, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Fprintf(ds.output, "TYPE: %s\n\n%s\n===\n", "request", string(data))
	if err != nil {
		return nil, err
	}

	res, err := ds.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	data, err = httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Fprintf(ds.output, "TYPE: %s\n\n%s\n===\n", "response", string(data))
	if err != nil {
		return nil, err
	}

	return res, nil
}
