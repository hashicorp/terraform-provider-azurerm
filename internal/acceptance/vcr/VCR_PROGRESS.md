# VCR Implementation Progress

## Overview
VCR (HTTP recording/playback) implementation for Terraform provider acceptance tests using go-vcr v4.

## Environment Variables
| Variable | Values | Description |
|----------|--------|-------------|
| `VCR_MODE` | `RECORD`, `REPLAY`, (empty) | Controls VCR mode. Empty = passthrough (real HTTP) |
| `VCR_PATH` | Absolute path | Directory for cassette files. Default: `.local/testdata/recordings` |
| `VCR_DEBUG` | `1`, `true`, `yes`, (empty) | Optional request/response debug logging for VCR transport |

## Cassette Path
```
VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings
```

## Test Command
```bash
# Recording mode
VCR_MODE=RECORD VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings \
  make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr'

# Replay mode
VCR_MODE=REPLAY VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings \
  make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr'

# Replay mode with transport debug logs
VCR_MODE=REPLAY VCR_DEBUG=1 VCR_PATH=~/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings \
  make acctests SERVICE='apimanagement' TESTARGS='-run=TestAccApiManagementLogger_testvcr'
```

## Test Case
- **Test**: `TestAccApiManagementLogger_testvcr`
- **File**: `internal/services/apimanagement/api_management_logger_resource_test.go`
- **Uses**: `data.ResourceTestWithVCR(t, r, []acceptance.TestStep{...})`

## Current Status Snapshot (2026-02-21)

- Replay is working for both:
  - `TestAccApiManagementLogger_testvcr`
  - `TestAccApiManagementLogger_basicEventHub`
- The previously observed `basicEventHub` replay mismatch (`200` in VCR debug vs check-path `404`) is no longer reproducing.
- Replay execution time has been reduced significantly and is now in the ~39s range for `TestAccApiManagementLogger_basicEventHub` in latest run.

## Completed Phases

### Phase 1: Passthrough Mode ✅
- VCR package created at `internal/acceptance/vcr/vcr.go`
- `GetHTTPClient(t)` returns HTTP client based on VCR_MODE
- Passthrough mode uses real HTTP requests (no recording)

### Phase 2: Record/Replay Mode ✅
- Recording mode saves HTTP interactions to YAML cassettes
- Replay mode reads from cassettes
- Cassette files stored at `{VCR_PATH}/{TestName}.yaml`
- Recorder only persists cassette/random artifacts on successful test completion
- On test failure/panic: recorder is intentionally not stopped and random pending state is cleared

### Phase 3: Deterministic Random Values ✅
- `RandTimeIntVCR(t,value)` - deterministic random integer for VCR tests
- `RandStringVCR(t,value)` - deterministic random string for VCR tests
- `IsVCRActive()` - returns true if VCR is RECORD or REPLAY
- Recording captures random values in memory and writes `{cassette}.random.json` only on successful completion
- Replay loads random values from `{cassette}.random.json` and reuses them

### Phase 4: Replay Validation ✅
- Baseline acceptance run (no VCR) completed
- Recording run completed and cassette/random file verified
- Replay run completed successfully for `TestAccApiManagementLogger_testvcr`
- Replay matching issue fixed by using a stable matcher that ignores volatile headers

### Phase 4b: Broader Replay Validation ✅
- `TestAccApiManagementLogger_basicEventHub` replay path now executes successfully.
- Transport-level replay is confirmed (microsecond response timings from cassette-backed requests).
- Remaining delays are mostly framework/poller orchestration overhead rather than real network latency.

### Phase 4c: Replay Performance Tuning ✅
- Added replay-aware handling for polling-related delays in VCR and poller paths.
- Applied replay-safe interval reductions in targeted APIM, resource-group, and SDK pollers.
- Revalidated repeatedly with replay logs (`api_replay` through `api_replay_6`) and observed progression:
  - ~393s → ~120s → ~93s → ~39s.

## Modified Files

### internal/acceptance/vcr/vcr.go
- `GetHTTPClient(t)` - VCR HTTP client with recorder
- Request matcher configured via `cassette.NewDefaultMatcher(...)` to ignore volatile request metadata
- Authorization and correlation request ID are removed from captured request headers
- `VCR_DEBUG` gated debug transport logging
- Replay transport overrides `Retry-After` to `0` in replay mode to avoid artificial poll delays

### internal/services/apimanagement/custompollers/api_management_poller.go
- Replay-aware poll interval handling for APIM delete/purge polling behavior

### internal/services/apimanagement/custompollers/api_management_api_poller.go
- Replay-aware poll interval handling for APIM API async polling behavior

### internal/services/apimanagement/api_management_resource.go
- Replay-aware initial poll intervals for APIM service delete/purge custom pollers

### internal/services/resource/resource_group_resource.go
- Replay-aware initial intervals for resource-group create/delete safety pollers

### internal/services/resource/custompollers/resource_group_prevent_delete_poller.go
- Replay-aware in-progress interval for nested-resource prevent-delete checks

### internal/services/resource/custompollers/resource_group_create_poller.go
- Replay-aware interval usage for resource-group create stabilization polling

### vendor/github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager/poller.go
- Replay-aware default polling interval for provisioning/delete poller creation paths

### vendor/github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager/poller_lro.go
- Replay-aware default interval for long-running operation pollers when `Retry-After` is absent


### internal/acceptance/data.go
- `BuildTestData(t, ...)` 

### internal/acceptance/testcase.go
- `ResourceTestWithVCR(t, r, steps)` - test runner with VCR HTTP client

### internal/provider/framework/factory_builder.go
- `ProtoV5ProviderFactoriesInitWithHTTPClient(httpClient)` - provider factory with custom HTTP client

### internal/provider/provider.go
- `AzureProviderWithHTTPClient(httpClient)` - provider schema with HTTP client injection

### internal/clients/builder.go
- `ClientBuilder.HttpClient` - field for custom HTTP client
- Passes HTTP client to ClientOptions

### internal/common/client_options.go
- `ClientOptions.HTTPClient` - field for custom HTTP client
- Sets HTTP client on SDK clients

### vendor/github.com/hashicorp/go-azure-sdk/sdk/client/client.go (patched)
- `SetHTTPClient(client)` - inject custom HTTP client
- `GetHTTPClient()` - retrieve HTTP client
- Modified `retryableClient()` to use injected client

## Architecture Flow
```
Test
  → BuildTestData(t, ...) 
  → data.ResourceTestWithVCR(t, r, steps)
      → vcr.GetHTTPClient(t) [creates recorder]
      → ProtoV5ProviderFactoriesInitWithHTTPClient(httpClient)
          → AzureProviderWithHTTPClient(httpClient)
              → ClientBuilder{HttpClient: httpClient}
                  → ClientOptions{HTTPClient: httpClient}
                      → SDK client.SetHTTPClient(httpClient)
                          → retryableClient uses httpClient
                              → VCR recorder intercepts HTTP
```

## Findings Summary

- Record mode remains expectedly slow (real Azure calls).
- Replay-mode latency is now largely non-network overhead:
  - Terraform plugin-testing command lifecycle overhead
  - Provider process reattach/start-stop overhead per CLI phase
  - Poller/wait loops that are independent of HTTP transport latency
- VCR interception is functioning correctly for the critical request paths.

## Next Steps

### Phase 5: Cleanup
- [ ] Expand sensitive-data sanitization beyond request headers (if needed)
- [ ] Normalize and centralize replay-aware poll interval helpers to reduce duplicated logic

### Phase 6: Integration
- [x] Document usage patterns for contributors
- [ ] Add contributor note on expected residual replay overhead (framework orchestration)

## Known Issues
- Cassettes may contain sensitive data (Bearer tokens, connection strings) - needs sanitization
- Small residual replay gaps (few seconds) can still appear due to framework/poller orchestration even when HTTP replay is near-instant

## Cassette Location
```
/Users/jitendragangwar/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings/TestAccApiManagementLogger_testvcr.yaml
```
