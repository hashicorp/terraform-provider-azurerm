// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"runtime"
	"time"
)

// TransportMode describes how the configured transport is being used.
type TransportMode string

const (
	// TransportModeDefault indicates that the transport is being used in the default mode, without recording or replaying traffic.
	TransportModeDefault TransportMode = ""

	// TransportModeVCRRecord indicates that the transport is recording live traffic.
	TransportModeVCRRecord TransportMode = "vcr_record"

	// TransportModeVCRReplay indicates that the transport is replaying previously recorded traffic.
	TransportModeVCRReplay TransportMode = "vcr_replay"

	// TransportModeVCRReplayWithNewEpisodes indicates that the transport is replaying previously recorded traffic, but will record new interactions if a replay miss occurs.
	TransportModeVCRReplayWithNewEpisodes TransportMode = "vcr_replay_with_new_episodes"
)

// GetDefaultHttpTransport returns a new default transport configured for SDK HTTP traffic.
func GetDefaultHttpTransport() http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			d := &net.Dialer{Resolver: &net.Resolver{}}
			return d.DialContext(ctx, network, addr)
		},
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
}
