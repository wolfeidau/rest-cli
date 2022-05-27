# rest-cli

Command line REST API client which signs requests using AWS [Signature Version 4 signing process](https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html) for use with [Amazon API Gateway](https://aws.amazon.com/api-gateway/).

[![Go Report Card](https://goreportcard.com/badge/github.com/wolfeidau/rest-cli)](https://goreportcard.com/report/github.com/wolfeidau/rest-cli)
[![Documentation](https://godoc.org/github.com/wolfeidau/rest-cli?status.svg)](https://godoc.org/github.com/wolfeidau/rest-cli)

# Usage

```
Usage: rest-cli <url>

A basic REST cli.

Arguments:
  <url>

Flags:
  -h, --help                     Show context-sensitive help.
      --version
      --verbose
      --dump
  -X, --method="GET"
      --service="execute-api"
  -d, --data=STRING
  -H, --headers=KEY=VALUE;...
```

# Examples

Send a `GET` request to an endpoint URL and pipe the output.

```
rest-cli https://example.com/customers
```

Send a `POST` request to an endpoint URL with a JSON body.

```
rest-cli -X POST -H 'Content-Type=application/json' -d '{"name":"AWS", "labels":[]}' https://example.com/customers
```

# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).
