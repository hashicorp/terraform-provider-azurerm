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

const (
	SubscriptionPlaceholder     = "00000000-0000-0000-0000-000000000000"
	SubscriptionPlaceholderAlt  = "00000000-0000-0000-0000-000000000001"
	SubscriptionPlaceholderAlt2 = "00000000-0000-0000-0000-000000000002"

	ModeRecordOnly            = recorder.ModeRecordOnly
	ModeReplayOnly            = recorder.ModeReplayOnly
	ModeReplayWithNewEpisodes = recorder.ModeReplayWithNewEpisodes
	ModeRecordOnce            = recorder.ModeRecordOnce
	ModePassthrough           = recorder.ModePassthrough

	vcrModeReplay = "replay"
	vcrModeRecord = "record"
)

var (
	recorders    = make(map[string]*recorder.Recorder)
	mu           sync.Mutex
	testDataPath = "vcrtestdata"

	// subscriptionRe matches common Azure subscription ID patterns in URLs and JSON
	subscriptionRe = regexp.MustCompile(`(?i)(/subscriptions/|subscriptionId=|subscription_id=|"subscriptionId":\s*")([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)

	// Where we have tests that use variable data from "random" sources, we need to catch this for VCR and make it predictable or we can't replay reliably.
	dynamicExtractors = map[string]*regexp.Regexp{
		// time.Now().String() for expirationTime
		"2045-01-01T00:00:00Z":                 regexp.MustCompile(`(?i)"expirationTime"\s*:\s*"([^"]+)"`),
		"11111111-1111-1111-1111-111111111111": regexp.MustCompile(`(?i)(?:/roleAssignments/|roleAssignmentId"?[=:]\s*"?)([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`),
	}
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

	redactSubscriptions := func(s string) string {
		return RedactSubscriptions(s, idReplacements)
	}

	redactHeaders := func(headers http.Header) {
		RedactHeaders(headers, idReplacements)
	}

	if r, exists := recorders[testName]; exists {
		return r, nil
	}

	// default to passthrough, just in case something unexpected is set
	mode := recorder.ModePassthrough
	vcrMode := os.Getenv("TC_TEST_VIA_VCR")
	switch vcrMode {
	case vcrModeRecord:
		mode = recorder.ModeRecordOnly
	case vcrModeReplay:
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
		// 1. Extract dynamic seeds from incoming Request
		seeds := make(map[string]string)
		for placeholder, ext := range dynamicExtractors {
			if match := ext.FindStringSubmatch(r.URL.String()); len(match) > 1 {
				seeds[placeholder] = match[1]
			}
		}

		var bodyStr string
		if r.Body != nil && r.Body != http.NoBody {
			if bodyBytes, err := io.ReadAll(r.Body); err == nil {
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // restore original
				bodyStr = string(bodyBytes)
				for placeholder, ext := range dynamicExtractors {
					if match := ext.FindStringSubmatch(bodyStr); len(match) > 1 {
						seeds[placeholder] = match[1]
					}
				}
			}
		}

		rCopy := r.Clone(r.Context())
		if bodyStr != "" {
			redactedBody := ScrubDynamicValues(redactSubscriptions(bodyStr), seeds)
			rCopy.Body = io.NopCloser(strings.NewReader(redactedBody))
			rCopy.ContentLength = int64(len(redactedBody))
		}

		// Normalise the incoming request before matching
		normalisedURL, err := url.Parse(ScrubDynamicValues(redactSubscriptions(r.URL.String()), seeds))
		if err == nil {
			rCopy.URL = normalisedURL
		}
		rCopy.RequestURI = ScrubDynamicValues(redactSubscriptions(rCopy.RequestURI), seeds)
		redactHeaders(rCopy.Header)

		// Cassette interaction already contains native placeholders because it was scrubbed at save
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
			// WrapTransport populates the sidecar in real-time as requests flow through,
			// so by the time this hook fires (at Stop()), all dynamic values for the session
			// are already present. We just read and apply them.
			seeds := readDynamicSidecar(testName)

			i.Request.URL = ScrubDynamicValues(redactSubscriptions(i.Request.URL), seeds)
			i.Request.RequestURI = ScrubDynamicValues(redactSubscriptions(i.Request.RequestURI), seeds)
			i.Request.Body = ScrubDynamicValues(redactSubscriptions(i.Request.Body), seeds)
			redactHeaders(i.Request.Headers)
			i.Response.Body = ScrubDynamicValues(redactSubscriptions(i.Response.Body), seeds)
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

// WrapTransport wraps a VCR Recorder to catch ErrInteractionNotFound and return an HTTP 418 response in tests.
// Prevents retry mechanisms from polling forever if a VCR interaction is missing.
// It also implements dynamic response rewriting to ensure timestamps stay matched across tests.
func WrapTransport(r *recorder.Recorder, testName string) http.RoundTripper {
	return &fallbackTransport{
		r:        r,
		testName: testName,
	}
}

func RedactHeaders(headers http.Header, idReplacements map[string]string) {
	for k, vals := range headers {
		for j, v := range vals {
			headers[k][j] = RedactSubscriptions(v, idReplacements)
		}
	}
}

// StopRecorder stops and removes the recorder from the map, saving it to disk.
func StopRecorder(testName string, failed bool) error {
	mu.Lock()
	defer mu.Unlock()

	if r, exists := recorders[testName]; exists {
		var err error
		saveOnFailure := os.Getenv("TC_VCR_SAVE_ON_FAILURE") != ""
		vcrMode := os.Getenv("TC_TEST_VIA_VCR")

		switch {
		case !failed, saveOnFailure, vcrMode != vcrModeRecord:
			if err = r.Stop(); err != nil {
				return fmt.Errorf("failed to stop recorder for %s: %v", testName, err)
			}

		default:
			// cassettePath := filepath.Join(testDataPath, testName) + ".yaml.gz" // switch these if using the gzip file handler
			cassettePath := filepath.Join(testDataPath, testName) + ".yaml"
			_ = os.Remove(cassettePath)
		}

		delete(recorders, testName)
	}

	// Always cleanup the dynamic sidecar
	cachePath := filepath.Join(testDataPath, testName+"_dynamic.json")
	_ = os.Remove(cachePath)

	return nil
}

func RedactSubscriptions(s string, idReplacements map[string]string) string {
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
