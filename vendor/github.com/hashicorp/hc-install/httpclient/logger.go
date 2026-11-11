// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package httpclient

import (
	"log"

	"github.com/hashicorp/go-retryablehttp"
)

func WithLogger(logger *log.Logger) Option {
	return func(c *retryablehttp.Client) {
		c.Logger = logger
	}
}
