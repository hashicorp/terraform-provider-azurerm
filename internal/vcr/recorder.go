// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vcr

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

var (
	recorders    = make(map[string]*recorder.Recorder)
	mu           sync.Mutex
	testDataPath = "vcrtestdata"
)

const subscriptionPlaceholder = "00000000-0000-0000-0000-000000000000"

// GetRecorder returns the shared recorder for a given test name, initialising it if necessary.
// it redacts sensitive information, such as subscriptionId and Authorization Headers and tailors the matcher to AzureRM
// requests.
func GetRecorder(testName string, subscriptionId string) (*recorder.Recorder, error) {
	if testName == "" {
		return nil, fmt.Errorf("testName must be provided to retrieve a recorder")
	}

	mu.Lock()
	defer mu.Unlock()

	subscriptionIDRe := regexp.MustCompile(subscriptionId)

	redactSubscriptions := func(s string) string {
		return subscriptionIDRe.ReplaceAllString(s, subscriptionPlaceholder)
	}

	if r, exists := recorders[testName]; exists {
		return r, nil
	}

	mode := recorder.ModeRecordOnce
	if os.Getenv("TC_TEST_VIA_VCR") == "record" {
		mode = recorder.ModeRecordOnly
	}

	// Cassette files are stored as gzip-compressed YAML (.yaml.gz).
	// The recorder library appends ".yaml" to the cassette name, so we suffix with ".gz"
	// to produce the final filename: vcrtestdata/<testName>.gz.yaml
	// For debugging/investigation purposes, switch the commented lines out and comment out the recorder.WithFS option in the recorder.New() below,
	// cassettePath := filepath.Join(testDataPath, testName+".gz")
	cassettePath := filepath.Join(testDataPath, testName)

	// ignore volatile per-request headers
	headerMatcher := cassette.NewDefaultMatcher(
		cassette.WithIgnoreAuthorization(),
		cassette.WithIgnoreUserAgent(),
		cassette.WithIgnoreHeaders(
			"X-Ms-Correlation-Request-Id",
			"X-Ms-Client-Request-Id",
			"X-Ms-Return-Client-Request-Id",
			"X-Ms-Routing-Request-Id",
			"X-Ms-Request-Id",
			"X-Msedge-Ref",
			"User-Agent",
		),
	)

	matcher := cassette.MatcherFunc(func(r *http.Request, i cassette.Request) bool {
		// Normalise subscription IDs in the incoming request URL before matching
		// so a real sub ID matches a cassette that has the placeholder.
		normalisedURL, err := url.Parse(redactSubscriptions(r.URL.String()))
		if err != nil {
			return false
		}
		rCopy := r.Clone(r.Context())
		rCopy.URL = normalisedURL
		// Also normalise in the cassette interaction copy
		iCopy := i
		iCopy.URL = redactSubscriptions(i.URL)
		return headerMatcher(rCopy, iCopy)
	})

	// Note: this needs to match the same struct from go-azure-sdk/sdk/client/client.go:retryableClient()
	defaultTransport := &http.Transport{
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

	r, err := recorder.New(cassettePath,
		recorder.WithMode(mode),
		recorder.WithSkipRequestLatency(true),
		recorder.WithRealTransport(defaultTransport),
		recorder.WithMatcher(matcher),
		// recorder.WithFS(&gzipFS{}),
		recorder.WithHook(func(i *cassette.Interaction) error {
			delete(i.Request.Headers, "Authorization")
			i.Request.Headers["Authorization"] = []string{"Bearer REDACTED"}
			return nil
		}, recorder.AfterCaptureHook),
		recorder.WithHook(func(i *cassette.Interaction) error {
			i.Request.URL = redactSubscriptions(i.Request.URL)
			i.Response.Body = redactSubscriptions(i.Response.Body)
			for k, vals := range i.Response.Headers {
				for j, v := range vals {
					i.Response.Headers[k][j] = redactSubscriptions(v)
				}
			}
			return nil
		}, recorder.AfterCaptureHook),

	)

	if err != nil {
		return nil, fmt.Errorf("failed to create recorder for %s: %v", testName, err)
	}

	recorders[testName] = r
	return r, nil
}

// StopRecorder stops and removes the recorder from the map, saving it to disk.
func StopRecorder(testName string) error {
	mu.Lock()
	defer mu.Unlock()

	if r, exists := recorders[testName]; exists {
		err := r.Stop()
		delete(recorders, testName)
		if err != nil {
			return fmt.Errorf("failed to stop recorder for %s: %v", testName, err)
		}
	}
	return nil
}
