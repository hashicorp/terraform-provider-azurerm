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

## Current Status Snapshot (2026-02-20)

- Accepted implementation remains focused on `TestAccApiManagementLogger_testvcr`.
- Replay for `TestAccApiManagementLogger_testvcr` is working.
- Replay for `TestAccApiManagementLogger_basicEventHub` is still failing.
- Failure signature observed in replay:
  - VCR debug shows a matching `GET .../loggers/...` with `200 OK`
  - Step check (`Check 1/5` in step `1/3`) fails with `404 ResourceGroupNotFound`
  - Error surfaces from `ExistsInAzure` path during check execution.

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

### Phase 4b: Broader Replay Validation ⚠️ (In Progress)
- `TestAccApiManagementLogger_basicEventHub` replay currently fails at step checks
- Additional diagnostics were added and used to capture full wrapped error chain
- Root-cause fix is still pending (candidate area: client path consistency between provider operations and check helpers)
- Resolve `basicEventHub` replay mismatch (`200` in VCR debug vs `404` in step check)
- Finalize/accept approach for VCR client lifecycle (singleton + shared client consistency)
- Re-validate in RECORD → REPLAY cycle for `TestAccApiManagementLogger_basicEventHub`

## Modified Files

### internal/acceptance/vcr/vcr.go
- `GetHTTPClient(t)` - VCR HTTP client with recorder
- Request matcher configured via `cassette.NewDefaultMatcher(...)` to ignore volatile request metadata
- Authorization and correlation request ID are removed from captured request headers
- `VCR_DEBUG` gated debug transport logging


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

## Next Steps

### Phase 5: Cleanup
- [ ] Expand sensitive-data sanitization beyond request headers (if needed)

### Phase 6: Integration
- [ ] Document usage patterns for contributors

## Known Issues
- Cassettes may contain sensitive data (Bearer tokens, connection strings) - needs sanitization
- Replay inconsistency remains for `TestAccApiManagementLogger_basicEventHub` despite successful `testvcr` replay

## Cassette Location
```
/Users/jitendragangwar/go/src/github.com/hashicorp/terraform-provider-azurerm/.local/testdata/recordings/TestAccApiManagementLogger_testvcr.yaml
```
