// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package vcr provides HTTP recording/playback for acceptance tests.
// This is a minimal POC implementation - currently a pass-through that makes real HTTP requests.
package vcr

import (
	"net/http"
	"testing"
)

// GetHTTPClient returns an HTTP client for acceptance tests.
// Currently returns a simple pass-through client for POC purposes.
// Future: Will support VCR recording/playback based on VCR_MODE environment variable.
func GetHTTPClient(t *testing.T) *http.Client {
	client := &http.Client{}
	return client
}
