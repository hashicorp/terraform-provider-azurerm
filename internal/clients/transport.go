// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"crypto/tls"
	"net"
	"net/http"
	"runtime"
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
// The settings mirror the go-azure-sdk default transport. Authorization is
// applied per request, so sharing the connection pool across clients and
// credentials is safe.
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
	MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
}
