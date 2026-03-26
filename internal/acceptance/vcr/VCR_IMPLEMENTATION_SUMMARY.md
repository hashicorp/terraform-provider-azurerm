# VCR Integration Summary for terraform-provider-azurerm

This document summarizes the VCR (HTTP recording/playback) integration work for faster acceptance test execution.

## Goal

Enable VCR-based testing using `go-vcr` to record HTTP interactions during test runs and replay them for faster subsequent test executions.

## Current Status

- ✅ Baseline acceptance run (no VCR) validated
- ✅ Recording mode (`VCR_MODE=RECORD`) validated
- ✅ Replay mode (`VCR_MODE=REPLAY`) validated for `TestAccApiManagementLogger_testvcr`
- ✅ Replay mode (`VCR_MODE=REPLAY`) validated for `TestAccApiManagementLogger_basicEventHub`
- ✅ Deterministic random values persisted/replayed via `{cassette}.random.json`
- ✅ Recorder persistence gated on test success only
- ✅ Failure/panic path skips recorder stop and clears pending random values
- ✅ Replay performance significantly improved for `TestAccApiManagementLogger_basicEventHub`

### Replay Performance Progression (basicEventHub)

- Baseline without VCR: ~349s
- Early replay attempt: ~393s
- After staged replay/poller optimizations: ~120s, then ~93s
- Latest replay run: ~39s

### Current Performance Findings

- Replay HTTP transport is working as expected (request/response replay in microseconds).
- Remaining wall-clock time is mostly from test harness and poller orchestration overhead rather than network latency.

## Current Architecture (Additional Param Approach)

The HTTPClient is passed as an additional parameter through the provider initialization chain.

### Data Flow

```
Test → vcr.GetHTTPClient(t) → ProtoV5ProviderFactoriesInitWithHTTPClient 
     → protoV5ProviderServerFactoryWithHTTPClient → AzureProviderWithHTTPClient 
     → azureProvider(httpClient) → buildClient(httpClient) → ClientBuilder.HttpClient 
     → ClientOptions.HTTPClient → SDK clients
```

## Modified Files

### 1. `internal/acceptance/vcr/vcr.go`
**Purpose:** VCR helper for HTTP recording/playback

**Current State:** Fully integrated VCR helper with:
- Mode handling (`RECORD`, `REPLAY`, passthrough)
- Absolute-path validation for `VCR_PATH`
- Sensitive request header stripping (`Authorization`, `X-Ms-Correlation-Request-Id`)
- Stable request matcher (ignores volatile metadata)
- Optional transport debug logging via `VCR_DEBUG`
- Deterministic random capture/replay (`.random.json`)
- Success-only persistence semantics for cassette/random artifacts
- Replay response header normalization for poller compatibility (`Retry-After: 0` in replay)

```go
const (
	EnvVCRMode  = "VCR_MODE"
	EnvVCRPath  = "VCR_PATH"
	EnvVCRDebug = "VCR_DEBUG"
	ModeRecord  = "RECORD"
	ModeReplay  = "REPLAY"
)

// GetHTTPClient returns nil in passthrough mode, otherwise a recorder-backed
// client configured for record or replay.
func GetHTTPClient(t *testing.T) *http.Client { /* ... */ }
```

### 2. `internal/acceptance/testcase.go`
**Purpose:** Test framework helpers

**Key Addition:** `ResourceTestWithVCR` method for opt-in VCR tests

```go
// ResourceTestWithVCR is an opt-in test method that uses VCR for HTTP recording/playback.
func (td TestData) ResourceTestWithVCR(t *testing.T, testResource types.TestResource, steps []TestStep) {	
	testCase := addStepsHelper(t, steps, td, testResource)
	td.runAcceptanceTestWithVCR(t, testCase)
}

// runAcceptanceTestWithVCR runs acceptance test with VCR HTTP client for recording/playback.
func (td TestData) runAcceptanceTestWithVCR(t *testing.T, testCase resource.TestCase) {
	testCase.ExternalProviders = td.externalProviders()
	httpClient := vcr.GetHTTPClient(t)
	testCase.ProtoV5ProviderFactories = framework.ProtoV5ProviderFactoriesInitWithHTTPClient(
		context.Background(), httpClient, "azurerm", "azurerm-alt")
	resource.ParallelTest(t, testCase)
}
```

### 3. `internal/provider/framework/factory_builder.go`
**Purpose:** Provider factory creation for tests

**Key Addition:** `ProtoV5ProviderFactoriesInitWithHTTPClient` function

```go
// ProtoV5ProviderFactoriesInitWithHTTPClient creates provider factories with a custom HTTP client.
func ProtoV5ProviderFactoriesInitWithHTTPClient(ctx context.Context, httpClient *http.Client, providerNames ...string) map[string]func() (tfprotov5.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov5.ProviderServer, error), len(providerNames))
	for _, name := range providerNames {
		factories[name] = func() (tfprotov5.ProviderServer, error) {
			providerServerFactory, _, err := protoV5ProviderServerFactoryWithHTTPClient(ctx, httpClient)
			if err != nil {
				return nil, err
			}
			return providerServerFactory(), nil
		}
	}
	return factories
}

func protoV5ProviderServerFactoryWithHTTPClient(ctx context.Context, httpClient *http.Client) (func() tfprotov5.ProviderServer, *schema.Provider, error) {
	v2Provider := provider.AzureProviderWithHTTPClient(httpClient)
	providers := []func() tfprotov5.ProviderServer{
		v2Provider.GRPCProvider,
		providerserver.NewProtocol5(NewFrameworkProvider(v2Provider)),
	}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		return nil, nil, err
	}
	return muxServer.ProviderServer, v2Provider, nil
}
```

### 4. `internal/provider/provider.go`
**Purpose:** Main provider implementation

**Key Addition:** `AzureProviderWithHTTPClient` function

```go
func AzureProvider() *schema.Provider {
	return azureProvider(false, nil)
}

// AzureProviderWithHTTPClient creates a provider with a custom HTTP client.
func AzureProviderWithHTTPClient(httpClient *http.Client) *schema.Provider {
	return azureProvider(false, httpClient)
}

func azureProvider(supportLegacyTestSuite bool, httpClient *http.Client) *schema.Provider {
	// ... existing code ...
	// httpClient is passed to buildClient
}
```

**In `buildClient` function:**
```go
func buildClient(ctx context.Context, p *schema.Provider, d *schema.ResourceData, httpClient *http.Client) (*clients.Client, diag.Diagnostics) {
	// ... existing code ...
	client, err := clients.Build(ctx, clients.ClientBuilder{
		// ... other fields ...
		HttpClient: httpClient,
	})
}
```

### 5. `internal/clients/builder.go`
**Purpose:** Client construction

**Key Addition:** `HttpClient` field in `ClientBuilder`

```go
type ClientBuilder struct {
	AuthConfig *auth.Credentials
	Features   features.UserFeatures
	// ... other fields ...
	HttpClient *http.Client  // <-- Added
}

func Build(ctx context.Context, builder ClientBuilder) (*Client, error) {
	// ... existing code ...
	o := &common.ClientOptions{
		// ... other fields ...
		HTTPClient: builder.HttpClient,  // <-- Passed to ClientOptions
	}
}
```

### 6. `internal/common/client_options.go`
**Purpose:** Common client options shared across services

**Key Addition:** `HTTPClient` field

```go
type ClientOptions struct {
	// ... other fields ...
	HTTPClient *http.Client
}

func (o ClientOptions) Configure(c client.BaseClient, authorizer auth.Authorizer) {
    // ... existing code ...
    c.SetHTTPClient(o.HTTPClient) // -> setting vcr http client in go-azure-sdk's baseclient
}
```

### 7. Replay Performance Optimizations (Additional)

The following replay-focused polling and wait adjustments were made:

- `internal/services/apimanagement/custompollers/api_management_poller.go`
- `internal/services/apimanagement/custompollers/api_management_api_poller.go`
- `internal/services/apimanagement/api_management_resource.go`
- `internal/services/resource/resource_group_resource.go`
- `internal/services/resource/custompollers/resource_group_prevent_delete_poller.go`
- `internal/services/resource/custompollers/resource_group_create_poller.go`
- `vendor/github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager/poller.go`
- `vendor/github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager/poller_lro.go`

These changes make polling intervals replay-aware to remove artificial 5s/10s/20s waits during cassette replay.


## How to Use VCR in Tests

To opt-in a test for VCR, change:
```go
// Before (normal test)
func TestAccApiManagementLogger_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{...})
}

// After (VCR-enabled test)
func TestAccApiManagementLogger_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_logger", "test")
	r := ApiManagementLoggerResource{}
	data.ResourceTestWithVCR(t, r, []acceptance.TestStep{...})  // <-- Use ResourceTestWithVCR
}
```


## Environment Variables

| Variable | Values | Description |
|----------|--------|-------------|
| `VCR_MODE` | `RECORD`, `REPLAY`, or empty | Controls VCR mode. Empty disables VCR. |
| `VCR_PATH` | absolute path | Directory for cassettes and random data files |
| `VCR_DEBUG` | `1`, `true`, `yes`, or empty | Enables VCR transport request/response debug logs |


## Test Execution

```bash
# Record mode - makes real API calls and saves responses
VCR_MODE=RECORD VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings \
	make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'

# Playback mode - uses recorded responses
VCR_MODE=REPLAY VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings \
	make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'

# Playback mode with debug tracing
VCR_MODE=REPLAY VCR_DEBUG=1 VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings \
	make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'

# Normal mode - VCR disabled, behaves like regular test
make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'
```

## Quickstart

Use these commands from repository root:

```bash
# 1) Normal execution (no VCR)
make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'

# 2) Record HTTP interactions
VCR_MODE=RECORD VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings \
	make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'

# 3) Replay from recorded cassette
VCR_MODE=REPLAY VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings \
	make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'
```

Expected artifacts after step 2:
- `{VCR_PATH}/TestAccApiManagementLogger_testvcr.yaml`
- `{VCR_PATH}/TestAccApiManagementLogger_testvcr.random.json`

Optional debug for replay diagnostics:

```bash
VCR_MODE=REPLAY VCR_DEBUG=1 VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings \
	make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr' TESTTIMEOUT='90m'
```

## Replay Matching Note

`go-vcr` default matching is strict and can fail replay when volatile headers/request metadata differ across runs. The implementation now uses `cassette.NewDefaultMatcher(...)` with ignored volatile fields.

## Current Open Considerations

- Replay is now functionally working for the targeted APIM logger tests.
- Residual gaps in replay logs can still appear due to framework lifecycle overhead (Terraform command phases, provider process orchestration) and not necessarily due to VCR transport.
- Further reductions are possible, but likely require deeper changes to acceptance framework execution strategy rather than VCR transport itself.

## Scope Note

This summary reflects accepted changes currently present in the branch. Experimental fixes that were not accepted are intentionally not treated as implemented state.

## File Structure

```
internal/
├── acceptance/
│   ├── testcase.go              # ResourceTestWithVCR, runAcceptanceTestWithVCR
│   └── vcr/
│       ├── vcr.go               # vcr recorder logic

├── clients/
│   ├── builder.go               # ClientBuilder.HttpClient
│   └── client.go                # Client struct
├── common/
│   └── client_options.go        # ClientOptions.HTTPClient
└── provider/
    ├── provider.go              # AzureProviderWithHTTPClient, buildClient
    └── framework/
        └── factory_builder.go   # ProtoV5ProviderFactoriesInitWithHTTPClient
```
