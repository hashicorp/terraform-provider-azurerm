# VCR Integration Summary for terraform-provider-azurerm

This document is the single current reference for VCR (HTTP recording/playback) support used by acceptance tests in this repository.

## Goal

Enable VCR-based testing using `go-vcr` so acceptance tests can record Azure HTTP interactions once and replay them in later runs.

## Current Status

- ✅ Baseline acceptance run (no VCR) validated
- ✅ Recording mode (`VCR_MODE=RECORD`) validated
- ✅ Replay mode (`VCR_MODE=REPLAY`) validated for `TestAccApiManagementLogger_testvcr`
- ✅ Replay mode (`VCR_MODE=REPLAY`) validated for `TestAccApiManagementLogger_basicEventHub`
- ✅ Deterministic random values persisted and replayed via `{cassette}.random.json`
- ✅ Recorder persistence is gated on test success only
- ✅ Failure/panic path skips recorder stop and clears pending random values
- ✅ Replay performance significantly improved for targeted APIM logger tests

### Replay Performance Progression (`basicEventHub`)

- Baseline without VCR: ~349s
- Early replay attempt: ~393s
- After staged replay optimizations: ~120s, then ~93s
- Latest replay run: ~39s

### Current Performance Findings

- Replay HTTP transport is working as expected
- Remaining wall-clock time is mostly test harness and poller orchestration overhead rather than network latency

## Environment Variables

| Variable | Values | Description |
|----------|--------|-------------|
| `VCR_MODE` | `RECORD`, `REPLAY`, or empty | Controls VCR mode. Matching is case-insensitive after trimming whitespace. Empty disables VCR. |
| `VCR_PATH` | absolute path, optional | Directory for cassettes and random data files. If unset, defaults to `{project-root}/.local/testdata/recordings`. |
| `VCR_DEBUG` | `1`, `true`, `yes`, or empty | Enables VCR transport request/response debug logs. |

## How VCR Works

The VCR HTTP client is passed through the provider initialization chain when VCR is active.

### Data Flow

```text
Test
  → data.ResourceTest(...)
      → vcr.IsVCRActive()
          → runAcceptanceTestWithVCR(...)
              → testclient.BuildWithVcr(t)
              → vcr.GetHTTPClient(t)
              → ProtoV5ProviderFactoriesInitWithHTTPClient(...)
                  → protoV5ProviderFactoriesInit(...)
                      → protoV5ProviderServerFactory(...)
                          → AzureProviderWithHTTPClient(...)
                              → buildClient(...)
                                  → ClientBuilder.HttpClient
                                      → ClientOptions.HTTPClient
                                          → SDK clients
```

## Main Implementation Details

### `internal/acceptance/vcr/vcr.go`

Current behavior:

- mode handling for `RECORD`, `REPLAY`, and passthrough
- normalized `VCR_MODE` parsing via trimming and uppercasing
- optional `VCR_PATH` with fallback to `{project-root}/.local/testdata/recordings`
- absolute-path validation after fallback resolution
- sensitive request header stripping for captured cassettes
- stable request matcher that ignores volatile request metadata
- optional transport debug logging via `VCR_DEBUG`
- deterministic random capture/replay using `.random.json`
- success-only persistence semantics for cassette/random artifacts
- replay response header normalization for poller compatibility via `Retry-After: 0`

### `internal/acceptance/testcase.go`

Current behavior:

- `ResourceTest(...)` is the standard acceptance entry point
- `ResourceTest(...)` conditionally routes into the VCR-aware execution path when `vcr.IsVCRActive()` is true
- `runAcceptanceTestWithVCR(...)` wires the recorder-backed HTTP client into provider factories
- `CheckDestroy` uses `testclient.BuildWithVcr(t)` in the VCR path so destroy checks use the same transport

### `internal/acceptance/testclient/client.go`

Current behavior:

- `BuildWithVcr(t)` creates a client using `vcr.GetHTTPClient(t)`
- this keeps acceptance-framework client creation aligned with the same recorder-backed HTTP transport

### `internal/provider/framework/factory_builder.go`

Current behavior:

- `ProtoV5ProviderFactoriesInitWithHTTPClient(...)` provides the VCR-aware entry point
- provider factory creation is deduplicated through shared helpers
- provider server construction is centralized in `protoV5ProviderServerFactory(...)`

### Provider/client plumbing

The recorder-backed HTTP client flows through:

- `internal/provider/provider.go` via `AzureProviderWithHTTPClient(...)`
- `internal/clients/builder.go` via `ClientBuilder.HttpClient`
- `internal/common/client_options.go` via `ClientOptions.HTTPClient`

### Shared SDK poller

`vendor/github.com/hashicorp/go-azure-sdk/sdk/client/pollers/poller.go` now centralizes replay wait reduction inside `PollUntilDone(...)`.

In replay mode:

- `retryDuration` is forced to `0`
- recorded poll responses are consumed immediately
- replay-specific delay reduction no longer depends on scattered per-poller overrides

## Usage

No test-code change is required for the standard acceptance flow.

Tests should continue to call `data.ResourceTest(...)`.

Example:

```go
func TestAccApiManagementLogger_basicEventHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		// ...
	})
}
```

When `VCR_MODE` is unset, the test runs normally. When `VCR_MODE=RECORD` or `VCR_MODE=REPLAY`, `ResourceTest(...)` automatically routes through the VCR-aware path.

## Commands

```bash
# Normal execution (no VCR)
make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'

# Record HTTP interactions using the default cassette directory
VCR_MODE=RECORD \
	make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'

# Replay from the recorded cassette using the default cassette directory
VCR_MODE=REPLAY \
	make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'

# Replay using an explicit cassette directory
VCR_MODE=REPLAY VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings \
	make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'

# Replay with transport debug tracing
VCR_MODE=REPLAY VCR_DEBUG=1 \
	make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'
```

Expected artifacts after a successful record run:

- `{VCR_PATH}/TestAccApiManagementLogger_testvcr.yaml`
- `{VCR_PATH}/TestAccApiManagementLogger_testvcr.random.json`

If `VCR_PATH` is omitted, artifacts are written under:

```text
{project-root}/.local/testdata/recordings
```

## Replay Matching Note

`go-vcr` default matching is strict and can fail replay when volatile request metadata differs across runs. The implementation uses `cassette.NewDefaultMatcher(...)` with ignored volatile fields.

## Poller Failure-Mode Note

The centralized replay optimization is functionally correct for successful cassette replay, but one behavior still matters during replay failures:

- `PollUntilDone(...)` can retry immediately in replay mode when a poller has `retryOnError=true`
- the currently identified in-repo usage is the resource-group nested-resource delete safety poller
- on replay mismatch or another retriable polling error, this can degrade into a tight retry loop until the polling context deadline is reached

This is not a happy-path correctness issue, but it is the main replay failure-mode to keep in mind when debugging cassette mismatch or poller replay failures.

## Current Open Considerations

- Replay is functionally working for the targeted APIM logger tests
- Residual replay gaps can still appear due to framework lifecycle overhead, provider process orchestration, and non-HTTP wait paths
- Further reductions are possible, but likely require deeper acceptance-framework changes rather than transport-level VCR changes

## File Structure

```text
internal/
├── acceptance/
│   ├── testcase.go              # ResourceTest, runAcceptanceTestWithVCR
│   ├── testclient/
│   │   └── client.go            # BuildWithVcr
│   └── vcr/
│       └── vcr.go               # VCR recorder logic
├── clients/
│   └── builder.go               # ClientBuilder.HttpClient
├── common/
│   └── client_options.go        # ClientOptions.HTTPClient
└── provider/
    ├── provider.go              # AzureProviderWithHTTPClient, buildClient
    └── framework/
        └── factory_builder.go   # Shared ProtoV5 provider factory helpers

vendor/
└── github.com/hashicorp/go-azure-sdk/sdk/client/pollers/
    └── poller.go                # Centralized replay-aware poll wait handling
```
