package rest

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Client struct {
	httpClient *http.Client
}

func New(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

func (cl *Client) DoRequest(method string, url string, headers map[string]string, body io.Reader, out io.Writer) error {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build request")
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	log.Debug().Fields(map[string]interface{}{
		"headers": req.Header,
	}).Msg("built headers")

	log.Info().Str("method", req.Method).Str("URL", req.URL.String()).Msg("built request")

	start := time.Now()

	res, err := cl.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send body: %w", err)
	}

	if res.StatusCode > 400 {
		return fmt.Errorf("bad response: %s", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	log.Info().Str("status", res.Status).Str("duration", time.Since(start).String()).Msg("request successful")

	fmt.Fprintln(out, string(data))

	return nil
}
