// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const maxDrainBytes = 1 << 20 // 1 MiB

// Options configures an *http.Transport created by New.
type Options struct {
	DisableProxy        bool
	DisableKeepAlives   bool
	IdleConnTimeout     time.Duration
	MaxIdleConnsPerHost int
}

// DefaultOptions returns the SDK's default transport configuration.
func DefaultOptions() Options {
	return Options{
		DisableProxy:        false,
		DisableKeepAlives:   false,
		IdleConnTimeout:     60 * time.Second,
		MaxIdleConnsPerHost: 16,
	}
}

// New returns a configured *http.Transport with connection pooling, timeouts,
// and HTTP/2 health-check pings enabled.
func New(opts Options) *http.Transport {
	var proxyFunc func(*http.Request) (*url.URL, error)
	if !opts.DisableProxy {
		proxyFunc = http.ProxyFromEnvironment
	}

	idleTimeout := opts.IdleConnTimeout
	if idleTimeout <= 0 {
		idleTimeout = 60 * time.Second
	}

	maxIdlePerHost := opts.MaxIdleConnsPerHost
	if maxIdlePerHost <= 0 {
		maxIdlePerHost = 16
	}

	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return &http.Transport{
		Proxy:             proxyFunc,
		DialContext:       dialer.DialContext,
		DisableKeepAlives: opts.DisableKeepAlives,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   maxIdlePerHost,
		IdleConnTimeout:       idleTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		HTTP2: &http.HTTP2Config{
			SendPingTimeout: 30 * time.Second,
			PingTimeout:     15 * time.Second,
		},
	}
}

// DrainAndClose drains up to 1 MiB from body and closes it so the underlying
// HTTP connection can be returned to the idle pool.
func DrainAndClose(body io.ReadCloser) {
	if body == nil || body == http.NoBody {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(body, maxDrainBytes))
	_ = body.Close()
}
