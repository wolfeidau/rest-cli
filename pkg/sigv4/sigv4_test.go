package sigv4_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/rest-cli/pkg/sigv4"
)

func TestTransport_RoundTrip(t *testing.T) {

	assert := require.New(t)

	type fields struct {
		awscfg      aws.Config
		region      string
		serviceName string
		body        io.Reader
		transport   http.RoundTripper
	}
	type args struct {
		req func(method, url string, body io.Reader) *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "sign and validate",
			fields: fields{
				awscfg: aws.Config{
					Credentials: creds(),
				},
				serviceName: "execute-api",
				region:      "us-east-1",
				body:        strings.NewReader("{}"),
				transport:   http.DefaultTransport,
			},
			args: args{
				req: request,
			},
			want: "200 OK",
		},
		{
			name: "sign and validate no body",
			fields: fields{
				awscfg: aws.Config{
					Credentials: creds(),
				},
				serviceName: "execute-api",
				region:      "us-east-1",
				body:        http.NoBody,
				transport:   http.DefaultTransport,
			},
			args: args{
				req: request,
			},
			want: "200 OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}))

			req := tt.args.req(http.MethodGet, server.URL, tt.fields.body)
			ts := sigv4.NewTransport(tt.fields.awscfg, tt.fields.serviceName, tt.fields.region, tt.fields.transport)
			got, err := ts.RoundTrip(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transport.RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Contains(req.Header.Get("Authorization"), "AWS4-HMAC-SHA256")

			if !reflect.DeepEqual(got.Status, tt.want) {
				t.Errorf("Transport.RoundTrip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func creds() aws.CredentialsProviderFunc {
	return func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     "test",
			SecretAccessKey: "test",
		}, nil
	}
}

func request(method, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)

	return req
}
