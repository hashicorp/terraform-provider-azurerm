// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package httpclient

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/hc-install/version"
)

func withUserAgent() Option {
	return func(c *retryablehttp.Client) {
		c.HTTPClient.Transport = &userAgentRoundTripper{
			inner:     c.HTTPClient.Transport,
			userAgent: fmt.Sprintf("hc-install/%s", version.Version()),
		}
	}
}

type userAgentRoundTripper struct {
	inner     http.RoundTripper
	userAgent string
}

func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", rt.userAgent)
	}
	return rt.inner.RoundTrip(req)
}
