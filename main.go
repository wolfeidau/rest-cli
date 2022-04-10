package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
		Method  string   `enum:"GET,POST,PATCH,PUT,DELETE" default:"GET"`
		Service string   `default:"execute-api"`
		URL     *url.URL `arg:""`
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

	req, err := http.NewRequest(flags.Method, flags.URL.String(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build request")
	}

	signer := v4.NewSigner()

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load creds")
	}

	err = signer.SignHTTP(ctx, creds, req, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", flags.Service, cfg.Region, time.Now())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load creds")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load creds")
	}

	if res.StatusCode != 200 {
		log.Error().Str("status", res.Status).Msg("request faild")

		return
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read body")
	}

	log.Info().Str("status", res.Status).Msg("request successful")

	fmt.Fprintln(os.Stdout, string(data))
}
