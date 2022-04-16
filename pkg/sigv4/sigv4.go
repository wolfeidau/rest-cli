package sigv4

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

const (
	// APIGWServiceName the service name for aws api gateway
	APIGWServiceName = "execute-api"
)

type Transport struct {
	signer      *v4.Signer
	awscfg      aws.Config
	region      string
	serviceName string
	transport   http.RoundTripper
}

func NewTransport(awscfg aws.Config, serviceName, region string, transport http.RoundTripper) *Transport {
	signer := v4.NewSigner()

	return &Transport{
		signer:      signer,
		awscfg:      awscfg,
		serviceName: serviceName,
		region:      region,
		transport:   transport,
	}
}

func (ts *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		r2  io.ReadCloser
		err error
	)

	req.Body, r2, err = copyBody(req.Body)
	if err != nil {
		return nil, err
	}

	payloadHash := generatePayloadHash(r2)

	creds, err := ts.awscfg.Credentials.Retrieve(req.Context())
	if err != nil {
		return nil, err
	}

	err = ts.signer.SignHTTP(req.Context(), creds, req, payloadHash, ts.serviceName, ts.region, time.Now())
	if err != nil {
		return nil, err
	}

	return ts.transport.RoundTrip(req)
}

func copyBody(b io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	if b == nil || b == http.NoBody {
		return http.NoBody, http.NoBody, nil
	}
	var (
		buf bytes.Buffer
		err error
	)
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func generatePayloadHash(body io.ReadCloser) string {
	h := sha256.New()
	_, _ = io.Copy(h, body)
	return hex.EncodeToString(h.Sum(nil))
}
