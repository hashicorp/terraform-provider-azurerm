// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package httpclient

import (
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

var discardLogger = log.New(io.Discard, "", 0)
var defaultOpts = []Option{
	WithLogger(discardLogger),
	withUserAgent(),
}

// New provides a pre-configured http.Client
// e.g. with relevant User-Agent header
func New(opts ...Option) *http.Client {
	rc := retryablehttp.NewClient()

	// process default options first
	for _, opt := range defaultOpts {
		opt(rc)
	}

	// process any other options
	for _, opt := range opts {
		opt(rc)
	}

	return rc.StandardClient()
}
