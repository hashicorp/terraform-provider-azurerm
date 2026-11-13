// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// sharedTransport is the process-wide http.RoundTripper used by every
// go-azure-sdk Resource Manager client built by this provider. Data-plane
// clients (common.ClientOptions.ConfigureDataPlane) do not use it.
//
// Without it, go-azure-sdk creates a new http.Transport for every request,
// so no connection is ever reused and each ARM call opens a new TCP+TLS
// connection. Behind source NAT with a fixed port budget (e.g. an Azure
// Standard Load Balancer) this exhausts the outbound ports. Embedding
// consumers that configure the provider repeatedly (e.g. Crossplane, once
// per reconcile) need the transport to outlive the configured client, so it
// is shared at package level.
//
// The settings start from the go-azure-sdk default transport, with two
// changes for a long-lived pool:
//   - MaxIdleConnsPerHost matches MaxIdleConns, because nearly all requests go
//     to one host. A lower limit closes connections after bursts on HTTP/1.1
//     paths (e.g. TLS-inspecting proxies), which brings back the churn.
//   - HTTP/2 health checks close a silently dropped connection after about
//     45 seconds. Otherwise every request multiplexed on it would hang until
//     its context deadline.
//
// Authorization is applied per request, so sharing the connection pool across
// clients and credentials is safe.
var sharedTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Resolver: &net.Resolver{},
	}).DialContext,
	TLSClientConfig: &tls.Config{
		MinVersion: tls.VersionTLS12,
	},
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	ForceAttemptHTTP2:     true,
	MaxIdleConnsPerHost:   100,
	HTTP2: &http.HTTP2Config{
		SendPingTimeout: 30 * time.Second,
		PingTimeout:     15 * time.Second,
	},
}
