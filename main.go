package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/rest-cli/pkg/rest"
	"github.com/wolfeidau/rest-cli/pkg/sigv4"
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

	httpClient := &http.Client{
		Transport: sigv4.NewTransport(cfg, flags.Service, cfg.Region, http.DefaultTransport),
	}

	var body io.Reader
	if len(flags.Data) > 0 {
		body = strings.NewReader(flags.Data)
	}

	restClient := rest.New(httpClient)

	err = restClient.DoRequest(flags.Method, flags.URL.String(), flags.Headers, body, os.Stdout)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to perform rest call")
	}
}
