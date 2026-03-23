// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package vcr provides HTTP recording/playback for acceptance tests.
// Uses go-vcr v4 for recording and replaying HTTP interactions.
package vcr

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	"time"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// VCR mode environment variable and values
const (
	EnvVCRMode  = "VCR_MODE"
	EnvVCRPath  = "VCR_PATH"
	EnvVCRDebug = "VCR_DEBUG"
	ModeRecord  = "RECORD"
	ModeReplay  = "REPLAY"
)

var (
	recorderSingletonMu sync.Mutex
	recorderSingletons  = make(map[string]*http.Client)
)

func isDebugEnabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(EnvVCRDebug)))
	return v == "1" || v == "true" || v == "yes"
}

// ensureRecordingsDir creates the recordings directory if it doesn't exist.
func ensureRecordingsDir(basePath string) error {
	return os.MkdirAll(basePath, 0755)
}

func validateAndGetBasePath(t *testing.T) string {
	basePath := os.Getenv(EnvVCRPath)
	if basePath == "" {
		projectRoot := findProjectRootPath(t)
		basePath = filepath.Join(projectRoot, ".local", "testdata", "recordings")
		t.Logf("%s not set, defaulting to: %s", EnvVCRPath, basePath)
	}

	if !filepath.IsAbs(basePath) {
		t.Fatalf("%s must be an absolute path, got: %s", EnvVCRPath, basePath)
	}
	return filepath.Clean(basePath)
}

func currentMode() string {
	return strings.ToUpper(strings.TrimSpace(os.Getenv(EnvVCRMode)))
}

// GetHTTPClient returns an HTTP client for acceptance tests.
// VCR mode is controlled by VCR_MODE environment variable:
//   - "RECORD": Record HTTP interactions to cassette file
//   - "REPLAY": Replay HTTP interactions from cassette file
//   - empty/other: Passthrough mode (make real HTTP requests, no recording)
func GetHTTPClient(t *testing.T) *http.Client {
	vcrMode := currentMode()
	if vcrMode == "" {
		return nil // Passthrough mode, return nil to indicate no VCR client
	}
	testName := t.Name()

	recorderSingletonMu.Lock()
	defer recorderSingletonMu.Unlock()
	if singleton, exists := recorderSingletons[testName]; exists {
		return singleton
	}
	baseDir := validateAndGetBasePath(t)
	cassettePath := getCassettePathFromBasePath(baseDir, t)
	debugEnabled := isDebugEnabled()

	var mode recorder.Mode

	switch vcrMode {
	case ModeRecord:
		mode = recorder.ModeRecordOnly
		if err := ensureRecordingsDir(baseDir); err != nil {
			t.Fatalf("Failed to create recordings directory: %v", err)
		}
	case ModeReplay:
		mode = recorder.ModeReplayOnly
	default:
		return nil
	}

	t.Logf("VCR mode: %s, cassette: %s", vcrMode, cassettePath)
	if debugEnabled {
		t.Logf("VCR debug enabled: base_dir=%s", baseDir)
	}

	tlsConfig := tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	httpTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			d := &net.Dialer{Resolver: &net.Resolver{}}
			return d.DialContext(ctx, network, addr)
		},
		TLSClientConfig:       &tlsConfig,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}

	sensitiveHeaderHook := func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		delete(i.Request.Headers, "X-Ms-Correlation-Request-Id")
		return nil
	}
	// Define how VCR will match requests to stored interactions.
	requestMatcher := cassette.NewDefaultMatcher(
		cassette.WithIgnoreUserAgent(),
		cassette.WithIgnoreAuthorization(),
		cassette.WithIgnoreHeaders("X-Ms-Correlation-Request-Id"),
	)

	rec, err := recorder.New(
		cassettePath,
		recorder.WithHook(sensitiveHeaderHook, recorder.AfterCaptureHook),
		recorder.WithMatcher(requestMatcher),
		recorder.WithMode(mode),
		recorder.WithRealTransport(httpTransport),
		recorder.WithSkipRequestLatency(true),
	)
	if err != nil {
		t.Fatalf("Failed to create VCR recorder in %s mode: %v", vcrMode, err)
	}

	// Ensure recorder is stopped when test completes (saves cassette in recording mode)
	t.Cleanup(func() {
		// Only persist artifacts (cassette/random values) for successful tests.
		// On failed tests (including panic), do not stop recorder so nothing is written.
		if t.Failed() {
			logFailureDiagnostics(t, vcrMode, cassettePath)
			clearTestRandomDataFromMap(t)
			t.Logf("VCR: test failed or panicked - skipping recorder stop and artifact persistence")
			return
		}

		if mode == recorder.ModeRecordOnly {
			if err := saveRandomDataForReplay(t); err != nil {
				t.Errorf("Failed to save VCR random data: %v", err)
				return
			}
		}
		if err := rec.Stop(); err != nil {
			t.Errorf("Failed to stop VCR recorder: %v", err)
		}
		recorderSingletonMu.Lock()
		delete(recorderSingletons, testName)
		recorderSingletonMu.Unlock()
	})

	// Return HTTP client with VCR transport
	client := rec.GetDefaultClient()
	baseTransport := client.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}
	if vcrMode == ModeReplay {
		baseTransport = &replayTransport{base: baseTransport}
	}
	if debugEnabled {
		baseTransport = &debugTransport{t: t, base: baseTransport, mode: vcrMode}
	}
	client.Transport = baseTransport

	recorderSingletons[testName] = client
	return client
}

type debugTransport struct {
	t    *testing.T
	base http.RoundTripper
	mode string
}

type replayTransport struct {
	base http.RoundTripper
}

func (r *replayTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := r.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp != nil && resp.Header != nil {
		resp.Header.Set("Retry-After", "0")
	}

	return resp, nil
}

func (d *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	d.t.Logf("%v VCR debug [%s] request: %s %.256s ", start, d.mode, req.Method, req.URL.String())

	resp, err := d.base.RoundTrip(req)
	if err != nil {
		d.t.Logf("%v VCR debug [%s] error after %s: %v", time.Now(), d.mode, time.Since(start), err)
		return nil, err
	}

	d.t.Logf("%v VCR debug [%s] response after %s: %.256s %s", time.Now(), d.mode, time.Since(start), req.URL.String(), resp.Status)
	return resp, nil
}

func getCassettePath(t *testing.T) string {
	return getCassettePathFromBasePath(validateAndGetBasePath(t), t)
}

func getCassettePathFromBasePath(basePath string, t *testing.T) string {
	testName := strings.ReplaceAll(t.Name(), "/", "_")
	cassettePath := filepath.Join(basePath, testName)

	return cassettePath
}

// IsVCRActive returns true if VCR is recording or replaying.
func IsVCRActive() bool {
	mode := currentMode()
	return mode == ModeRecord || mode == ModeReplay
}

// testRandomData stores random values for a test
type testRandomData struct {
	RandomInteger int    `json:"random_integer"`
	RandomString  string `json:"random_string"`
}

// cachedRandomData stores loaded random data per test to avoid repeated file reads
var cachedRandomData = make(map[string]*testRandomData)

// pendingRandomData stores random data captured during recording and persisted on test success.
var pendingRandomData = make(map[string]*testRandomData)
var randomDataMu sync.Mutex

// getRandomDataPath returns the path to the random data file for a test
func getRandomDataPath(t *testing.T) string {
	cassettePath := getCassettePath(t)
	return cassettePath + ".random.json"
}

// loadRandomData loads random values from a JSON file during replay
func loadRandomData(t *testing.T) (*testRandomData, error) {
	randomDataMu.Lock()
	defer randomDataMu.Unlock()

	// Check cache first
	if data, ok := cachedRandomData[t.Name()]; ok {
		return data, nil
	}

	path := getRandomDataPath(t)

	jsonData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read random data file %s: %w", path, err)
	}

	var data testRandomData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal random data: %w", err)
	}

	// Cache the loaded data
	cachedRandomData[t.Name()] = &data
	t.Logf("VCR: Loaded random data from %s (RandomInteger: %d, RandomString: %s)", path, data.RandomInteger, data.RandomString)
	return &data, nil
}

func loadRandomDataOrFatal(t *testing.T, valueType string) *testRandomData {
	data, err := loadRandomData(t)
	if err != nil {
		t.Fatalf("VCR replay failed to load random %s: %v", valueType, err)
	}
	return data
}

func storeRandomIntegerForReplay(t *testing.T, value int) {
	randomDataMu.Lock()
	defer randomDataMu.Unlock()

	data, ok := pendingRandomData[t.Name()]
	if !ok {
		data = &testRandomData{}
		pendingRandomData[t.Name()] = data
	}
	data.RandomInteger = value
}

func storeRandomStringForReplay(t *testing.T, value string) {
	randomDataMu.Lock()
	defer randomDataMu.Unlock()

	data, ok := pendingRandomData[t.Name()]
	if !ok {
		data = &testRandomData{}
		pendingRandomData[t.Name()] = data
	}
	data.RandomString = value
}

func saveRandomDataForReplay(t *testing.T) error {
	randomDataMu.Lock()
	data, ok := pendingRandomData[t.Name()]
	defer randomDataMu.Unlock()
	if !ok {
		return nil
	}
	delete(pendingRandomData, t.Name())

	path := getRandomDataPath(t)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal random data: %w", err)
	}
	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write random data file: %w", err)
	}
	t.Logf("VCR: Saved random data to %s (RandomInteger: %d, RandomString: %s)", path, data.RandomInteger, data.RandomString)
	return nil
}

func clearTestRandomDataFromMap(t *testing.T) {
	randomDataMu.Lock()
	defer randomDataMu.Unlock()
	delete(pendingRandomData, t.Name())
}

func logFailureDiagnostics(t *testing.T, vcrMode string, cassettePath string) {
	t.Helper()

	randomPath := getRandomDataPath(t)
	cassetteYamlPath := cassettePath + ".yaml"

	t.Logf("VCR diagnostics: mode=%s test=%s", vcrMode, t.Name())
	t.Logf("VCR diagnostics: cassette_base=%s cassette_yaml=%s random_data=%s", cassettePath, cassetteYamlPath, randomPath)

	if info, err := os.Stat(cassetteYamlPath); err == nil {
		t.Logf("VCR diagnostics: cassette exists (size=%d bytes, modified=%s)", info.Size(), info.ModTime().Format(time.RFC3339))
	} else {
		t.Logf("VCR diagnostics: cassette not present or unreadable: %v", err)
	}

	if info, err := os.Stat(randomPath); err == nil {
		t.Logf("VCR diagnostics: random data exists (size=%d bytes, modified=%s)", info.Size(), info.ModTime().Format(time.RFC3339))
	} else {
		t.Logf("VCR diagnostics: random data not present or unreadable: %v", err)
	}

	randomDataMu.Lock()
	if pending, ok := pendingRandomData[t.Name()]; ok {
		t.Logf("VCR diagnostics: pending random data before cleanup: random_integer=%d random_string=%s", pending.RandomInteger, pending.RandomString)
	} else {
		t.Logf("VCR diagnostics: pending random data before cleanup: none")
	}
	randomDataMu.Unlock()

	t.Logf("VCR diagnostics stack:\n%s", string(debug.Stack()))
}

// RandTimeIntVCR returns a random integer for VCR tests.
// In RECORD mode: returns the provided value (will be saved on test success)
// In REPLAY mode: returns the captured value from recording
func RandTimeIntVCR(t *testing.T, value int) int {
	switch currentMode() {
	case ModeRecord:
		// Capture for persistence on successful test completion.
		storeRandomIntegerForReplay(t, value)
		return value
	case ModeReplay:
		return loadRandomDataOrFatal(t, "integer").RandomInteger
	default:
		return value
	}
}

// RandStringVCR returns a random string for VCR tests.
// In RECORD mode: returns the provided value (will be saved on test success)
// In REPLAY mode: returns the captured value from recording
func RandStringVCR(t *testing.T, value string) string {
	switch currentMode() {
	case ModeRecord:
		// Capture for persistence on successful test completion.
		storeRandomStringForReplay(t, value)
		return value
	case ModeReplay:
		return loadRandomDataOrFatal(t, "string").RandomString
	default:
		return value
	}
}

func findProjectRootPath(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	current := filepath.Clean(wd)
	for {
		goModPath := filepath.Join(current, "go.mod")
		if info, statErr := os.Stat(goModPath); statErr == nil && !info.IsDir() {
			return current
		}

		parent := filepath.Dir(current)
		if parent == current {
			return filepath.Clean(wd)
		}
		current = parent
	}
}
