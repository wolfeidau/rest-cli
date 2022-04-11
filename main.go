package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	version = "unknown"

	flags struct {
		Verbose bool
		Version kong.VersionFlag
		Method  string            `enum:"GET,POST,PATCH,PUT,DELETE" default:"GET" short:"X"`
		Service string            `default:"execute-api"`
		Data    string            `short:"d"`
		Headers map[string]string `short:"H"`
		URL     *url.URL          `arg:""`
	}
)

func main() {

	kong.Parse(&flags,
		kong.Vars{"version": version},
		kong.Name("rest-cli"),
		kong.Description("A basic REST cli."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
			Summary: true,
		}))

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if flags.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	var body io.Reader
	if len(flags.Data) > 0 {
		body = strings.NewReader(flags.Data)
	}

	req, err := http.NewRequest(flags.Method, flags.URL.String(), body)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build request")
	}

	for k, v := range flags.Headers {
		req.Header.Add(k, v)
	}

	signer := v4.NewSigner()

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load creds")
	}

	payloadHash := generatePayloadHash(flags.Data)

	err = signer.SignHTTP(ctx, creds, req, payloadHash, flags.Service, cfg.Region, time.Now())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load creds")
	}

	if flags.Verbose {
		log.Debug().Str("method", req.Method).Str("URL", req.URL.String()).Fields(map[string]interface{}{
			"headers": req.Header,
		}).Msg("built request")
	} else {
		log.Info().Str("method", req.Method).Str("URL", req.URL.String()).Msg("built request")
	}

	start := time.Now()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load creds")
	}

	if res.StatusCode > 400 {
		log.Error().Str("status", res.Status).Msg("request failed")

		return
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read body")
	}

	log.Info().Str("status", res.Status).Str("duration", time.Since(start).String()).Msg("request successful")

	fmt.Fprintln(os.Stdout, string(data))
}

func generatePayloadHash(body string) string {
	reader := strings.NewReader(body)

	h := sha256.New()
	_, _ = io.Copy(h, reader)
	return hex.EncodeToString(h.Sum(nil))
}