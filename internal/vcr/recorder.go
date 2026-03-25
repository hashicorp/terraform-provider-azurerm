// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vcr

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
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

const (
	SubscriptionPlaceholder     = "00000000-0000-0000-0000-000000000000"
	SubscriptionPlaceholderAlt  = "00000000-0000-0000-0000-000000000001"
	SubscriptionPlaceholderAlt2 = "00000000-0000-0000-0000-000000000002"
)

// GetRecorder returns the shared recorder for a given test name, initialising it if necessary.
// it redacts sensitive information, such as SubscriptionID and Authorization Headers and tailors the matcher to AzureRM
// requests.
func GetRecorder(testName string, subscriptionId string) (*recorder.Recorder, error) {
	if testName == "" {
		return nil, errors.New("testName must be provided to retrieve a recorder")
	}

	mu.Lock()
	defer mu.Unlock()

	primary := os.Getenv("ARM_SUBSCRIPTION_ID")
	alt := os.Getenv("ARM_SUBSCRIPTION_ID_ALT")
	alt2 := os.Getenv("ARM_SUBSCRIPTION_ID_ALT2")

	// Map real ID to specific placeholder
	idReplacements := make(map[string]string)
	if subscriptionId != "" && subscriptionId != primary {
		idReplacements[subscriptionId] = SubscriptionPlaceholder
	}
	if primary != "" {
		idReplacements[primary] = SubscriptionPlaceholder
	}
	if alt != "" {
		idReplacements[alt] = SubscriptionPlaceholderAlt
	}
	if alt2 != "" {
		idReplacements[alt2] = SubscriptionPlaceholderAlt2
	}

	// subscriptionRe matches common Azure subscription ID patterns in URLs and JSON
	subscriptionRe := regexp.MustCompile(`(?i)(/subscriptions/|subscriptionId=|subscription_id=|"subscriptionId":\s*")([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)

	redactSubscriptions := func(s string) string {
		// 1. Redact specific known IDs with their mapped placeholder (handles headers and other standalone occurrences)
		for id, placeholder := range idReplacements {
			re := regexp.MustCompile("(?i)" + id)
			s = re.ReplaceAllString(s, placeholder)
		}
		// 2. Catch any other subscription IDs in standard patterns
		return subscriptionRe.ReplaceAllStringFunc(s, func(match string) string {
			groups := subscriptionRe.FindAllStringSubmatch(match, -1)
			if len(groups) > 0 && len(groups[0]) > 2 {
				// Prevent double-matching placeholders if already replaced
				if groups[0][2] == SubscriptionPlaceholder || groups[0][2] == SubscriptionPlaceholderAlt || groups[0][2] == SubscriptionPlaceholderAlt2 {
					return match
				}
				return groups[0][1] + SubscriptionPlaceholder
			}
			return match
		})
	}

	redactHeaders := func(headers http.Header) {
		for k, vals := range headers {
			for j, v := range vals {
				headers[k][j] = redactSubscriptions(v)
			}
		}
	}

	if r, exists := recorders[testName]; exists {
		return r, nil
	}

	// default to passthrough, just in case something unexpected is set
	mode := recorder.ModePassthrough
	vcrMode := os.Getenv("TC_TEST_VIA_VCR")
	switch vcrMode {
	case "record":
		mode = recorder.ModeRecordOnly
	case "replay", "true":
		mode = recorder.ModeReplayOnly
	}

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
		// Normalise subscription IDs in the incoming request before matching
		normalisedURL, err := url.Parse(redactSubscriptions(r.URL.String()))
		if err != nil {
			return false
		}
		rCopy := r.Clone(r.Context())
		rCopy.URL = normalisedURL
		rCopy.RequestURI = redactSubscriptions(rCopy.RequestURI)
		redactHeaders(rCopy.Header)

		// Redact Body in the incoming request so body matching succeeds
		if r.Body != nil && r.Body != http.NoBody {
			if bodyBytes, err := io.ReadAll(r.Body); err == nil {
				// Restore original body for proper processing downstream
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				redactedBody := redactSubscriptions(string(bodyBytes))
				rCopy.Body = io.NopCloser(strings.NewReader(redactedBody))
				rCopy.ContentLength = int64(len(redactedBody))
			}
		}

		// Also normalise in the cassette interaction copy
		iCopy := i
		iCopy.URL = redactSubscriptions(i.URL)
		iCopy.RequestURI = redactSubscriptions(i.RequestURI)
		iCopy.Body = redactSubscriptions(i.Body)
		redactHeaders(i.Headers)

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
		// recorder.WithFS(&GzipFS{}),
		recorder.WithHook(func(i *cassette.Interaction) error {
			delete(i.Request.Headers, "Authorization")
			i.Request.Headers["Authorization"] = []string{"Bearer REDACTED"}
			return nil
		}, recorder.BeforeSaveHook),
		recorder.WithHook(func(i *cassette.Interaction) error {
			i.Request.URL = redactSubscriptions(i.Request.URL)
			i.Request.RequestURI = redactSubscriptions(i.Request.RequestURI)
			i.Request.Body = redactSubscriptions(i.Request.Body)
			redactHeaders(i.Request.Headers)
			i.Response.Body = redactSubscriptions(i.Response.Body)
			redactHeaders(i.Response.Headers)
			return nil
		}, recorder.BeforeSaveHook),
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
