# PowerShell Test Suite

Test suite for the PowerShell installer (`install-copilot-setup.ps1`).

## Test Files

### `run-all-tests.ps1`
Comprehensive test runner for PowerShell installer functionality.

### `cleanup-test-env.ps1`
Test environment cleanup utilities.

### `TestCases/`
Individual test case files:
- `test-cleanup-edge-cases.ps1` - Cleanup edge case scenarios
- `test-cleanup-logic.ps1` - Cleanup logic validation
- `test-cleanup-scenarios.ps1` - Various cleanup scenarios  
- `test-edge-cases.ps1` - General edge case testing
- `test-fresh-install.ps1` - Fresh installation testing
- `test-regex-patterns.ps1` - Regex pattern validation

### `MockScenarios/`
Mock test scenarios and corrupted state simulations.

## Usage

```powershell
# From the test directory
cd .github\AILocalInstallTest\powershell
.\run-all-tests.ps1

# Run specific test cases
.\TestCases\test-fresh-install.ps1

# Clean up test environment
.\cleanup-test-env.ps1
```

## Requirements

- PowerShell 5.1+ (Windows PowerShell) or PowerShell 7+ (PowerShell Core)
- Access to `..\..\AILocalInstall\` directory
- Administrative privileges may be required for some tests

## Test Results

- ✅ **PASS**: Component works correctly
- ❌ **FAIL**: Component has issues
- ⚠️ **WARNING**: Potential issues detected
- ℹ️ **INFO**: Additional information

## Adding Tests

To add new PowerShell-specific tests:

1. Create new test file in `TestCases/` directory
2. Follow naming convention: `test-feature-name.ps1`
3. Use consistent output format for results
4. Update this README with test description
