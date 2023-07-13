// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"context"
	"crypto/tls"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

var (
	// Client is the HTTP client used for sending authentication requests and obtaining tokens
	Client HTTPClient

	// MetadataClient is the HTTP client used for obtaining tokens from the Instance Metadata Service
	MetadataClient HTTPClient
)

func init() {
	Client = httpClient(defaultHttpClientParams())
	MetadataClient = httpClient(httpClientParams{
		instanceMetadataService: true,

		retryWaitMin:  2 * time.Second,
		retryWaitMax:  60 * time.Second,
		retryMaxCount: 5,
		useProxy:      false,
	})
}

type httpClientParams struct {
	instanceMetadataService bool

	retryWaitMin  time.Duration
	retryWaitMax  time.Duration
	retryMaxCount int
	useProxy      bool
}

func defaultHttpClientParams() httpClientParams {
	return httpClientParams{
		instanceMetadataService: false,

		retryWaitMin:  1 * time.Second,
		retryWaitMax:  30 * time.Second,
		retryMaxCount: 8,
		useProxy:      true,
	}
}

// httpClient returns a shimmed retryablehttp Client, with custom backoff and
// retry settings which can be customized per instance as needed.
func httpClient(params httpClientParams) *http.Client {
	r := retryablehttp.NewClient()

	r.Logger = log.Default()

	r.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		// note: min and max contain the values of r.RetryWaitMin and r.RetryWaitMax

		if resp != nil {
			if params.instanceMetadataService {
				// IMDS uses inappropriate 410 status to indicate a rebooting-like state, retry after 70 seconds
				// See https://learn.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/how-to-use-vm-token#retry-guidance
				if resp.StatusCode == http.StatusGone {
					return 70 * time.Second
				}
			}

			// Always look for Retry-After header, regardless of HTTP status
			if s, ok := resp.Header["Retry-After"]; ok {
				if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
					return time.Second * time.Duration(sleep)
				}
			}
		}

		// Exponential backoff when Retry-After header not provided, e.g. IMDS
		mult := math.Pow(2, float64(attemptNum)) * float64(min)
		sleep := time.Duration(mult)
		if float64(sleep) != mult || sleep > max {
			sleep = max
		}
		return sleep
	}

	var proxyFunc func(*http.Request) (*url.URL, error)
	if params.useProxy {
		proxyFunc = http.ProxyFromEnvironment
	}

	r.RetryWaitMin = params.retryWaitMin
	r.RetryWaitMax = params.retryWaitMax

	tlsConfig := tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	r.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: proxyFunc,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := &net.Dialer{Resolver: &net.Resolver{}}
				return d.DialContext(ctx, network, addr)
			},
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSClientConfig:       &tlsConfig,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     true,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}

	return r.StandardClient()
}
