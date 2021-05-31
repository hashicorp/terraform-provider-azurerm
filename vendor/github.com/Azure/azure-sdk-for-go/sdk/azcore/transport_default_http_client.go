// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"crypto/tls"
	"net/http"
)

var defaultHTTPClient *http.Client

func init() {
	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	defaultTransport.TLSClientConfig.MinVersion = tls.VersionTLS12
	defaultHTTPClient = &http.Client{
		Transport: defaultTransport,
	}
}
