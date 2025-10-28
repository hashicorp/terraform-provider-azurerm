---
applyTo: "internal/**/*.go"
description: Troubleshooting decision trees and diagnostic patterns for the Terraform AzureRM provider including common issues, debugging workflows, and resolution strategies.
---

# 🔧 Troubleshooting Decision Trees

Troubleshooting decision trees and diagnostic patterns for the Terraform AzureRM provider including common issues, debugging workflows, and resolution strategies.

**Quick navigation:** [🚨 Common Issues](#🚨-common-issues) | [🔍 Debugging Workflows](#🔍-debugging-workflows) | [⚡ Quick Fixes](#⚡-quick-fixes) | [🏗️ Development Troubleshooting](#🏗️-development-troubleshooting)

## 🚨 Common Issues

### Azure API Rate Limiting

**Symptoms:**
- HTTP 429 errors in logs
- Intermittent failures during resource operations
- Slow resource creation/update cycles

**Decision Tree:**
```text
API Rate Limiting Detected
├─ Check subscription limits
│  ├─ Review Azure portal quotas
│  ├─ Verify service tier limits
│  └─ Consider subscription upgrade
├─ Implement retry logic
│  ├─ Use exponential backoff
│  ├─ Add jitter to reduce thundering herd
│  └─ Set maximum retry limits
└─ Optimize API calls
   ├─ Batch operations where possible
   ├─ Cache frequently accessed data
   └─ Reduce unnecessary API calls
```

**Resolution Pattern:**
```go
// Implement proper retry with exponential backoff
func retryWithBackoff(operation func() error) error {
    backoff := time.Second
    maxRetries := 5

    for i := 0; i < maxRetries; i++ {
        err := operation()
        if err == nil {
            return nil
        }

        if !isRetryableError(err) {
            return err
        }

        time.Sleep(backoff)
        backoff *= 2
        if backoff > 30*time.Second {
            backoff = 30*time.Second
        }
    }

    return fmt.Errorf("operation failed after %d retries", maxRetries)
}
```

### Resource State Drift

**Symptoms:**
- Terraform shows unexpected diffs on plan
- Resources appear modified outside Terraform
- Import operations fail with state mismatches

**Decision Tree:**
```text
State Drift Detected
├─ Identify drift source
│  ├─ Manual Azure portal changes
│  ├─ Other automation tools
│  ├─ Azure service auto-scaling
│  └─ Provider version differences
├─ Resolve drift
│  ├─ Update Terraform configuration to match
│  ├─ Import resources to sync state
│  ├─ Apply changes to restore desired state
│  └─ Use refresh-only plan to update state
└─ Prevent future drift
   ├─ Implement Azure Policy controls
   ├─ Use resource locks where appropriate
   ├─ Establish change management processes
   └─ Monitor for unauthorized changes
```

### Authentication and Authorization Issues

**Symptoms:**
- HTTP 401/403 errors
- "Principal does not have access" errors
- Authentication timeouts

**Decision Tree:**
```text
Authentication Issue
├─ Verify credentials
│  ├─ Check environment variables
│  ├─ Validate service principal
│  ├─ Confirm tenant/subscription IDs
│  └─ Test credential expiration
├─ Check permissions
│  ├─ Review Azure RBAC assignments
│  ├─ Verify resource-level permissions
│  ├─ Check API permissions for service principal
│  └─ Validate subscription access
└─ Test authentication
   ├─ Use Azure CLI for validation
   ├─ Test with minimal permissions
   ├─ Verify network connectivity
   └─ Check for conditional access policies
```

## 🔍 Debugging Workflows

### Step-by-Step Resource Debugging

**1. Information Gathering**
```bash
# Check Terraform version and provider version
terraform version

# Review resource configuration
terraform show -json | jq '.values.root_module.resources[] | select(.address == "azurerm_resource.example")'

# Check current state
terraform state show azurerm_resource.example
```

**2. Azure SDK Debugging**
```bash
# Enable detailed logging
$env:TF_LOG = "DEBUG"
$env:ARM_LOG_LEVEL = "DEBUG"

# Run targeted operation
terraform plan -target=azurerm_resource.example
```

**3. API Level Debugging**
```bash
# Use Azure CLI to test API directly
az rest --method GET --url "https://management.azure.com/subscriptions/{subscription-id}/resourceGroups/{rg}/providers/Microsoft.Service/resources/{name}?api-version=2023-01-01"
```

### Network and Connectivity Issues

**Debugging Pattern:**
```text
Connectivity Issue
├─ Test basic connectivity
│  ├─ Check internet connection
│  ├─ Verify DNS resolution
│  ├─ Test Azure endpoints
│  └─ Check proxy/firewall settings
├─ Azure-specific tests
│  ├─ Test authentication endpoint
│  ├─ Verify Azure API endpoints
│  ├─ Check service-specific endpoints
│  └─ Test from different networks
└─ Provider-specific debugging
   ├─ Enable TF_LOG=DEBUG
   ├─ Check HTTP response codes
   ├─ Review timeout settings
   └─ Test with reduced concurrency
```

## ⚡ Quick Fixes

### Common Error Resolution

**"Resource already exists" during creation:**
```bash
# Import existing resource
terraform import azurerm_resource.example /subscriptions/.../resourceGroups/.../providers/Microsoft.Service/resources/name

# Or force replacement
terraform apply -replace=azurerm_resource.example
```

**"Resource not found" during read:**
```bash
# Refresh state to detect deletion
terraform refresh

# Remove from state if manually deleted
terraform state rm azurerm_resource.example
```

**Schema validation errors:**
```hcl
# Check for deprecated arguments
# Review provider upgrade guides
# Validate argument types and values
```

### Performance Optimization

**Slow plan/apply operations:**
```bash
# Reduce parallelism
terraform plan -parallelism=1

# Target specific resources
terraform plan -target=azurerm_resource.example

# Use partial configuration
terraform plan -var-file=minimal.tfvars
```

## 🏗️ Development Troubleshooting

### Provider Development Issues

**Build Failures:**
```bash
# Check Go version compatibility
go version

# Update dependencies
go mod tidy

# Run specific tests
go test -v ./internal/services/servicename -run TestAccResourceName_basic
```

**Test Failures:**
```bash
# Run with detailed output
TF_ACC=1 go test -v ./internal/services/servicename -run TestAccResourceName_basic -timeout 60m

# Check for resource cleanup issues
# Review Azure credentials and permissions
# Verify test resource naming patterns
```

**Debugging Test Issues:**
```go
// Add debug logging to tests
t.Logf("Testing configuration: %s", config)

// Use acceptance.BuildTestData for consistent naming
data := acceptance.BuildTestData(t, "azurerm_resource", "test")

// Check for test isolation issues
// Verify resource group cleanup
// Review parallel test execution
```

### CustomizeDiff Debugging

**Validation Logic Issues:**
```go
// Add logging to CustomizeDiff functions
func validateConfiguration(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
    log.Printf("[DEBUG] CustomizeDiff: validating configuration")

    // Test specific field combinations
    enabled := diff.Get("enabled").(bool)
    config := diff.Get("configuration").([]interface{})

    log.Printf("[DEBUG] enabled: %t, config length: %d", enabled, len(config))

    if enabled && len(config) == 0 {
        return fmt.Errorf("`configuration` is required when `enabled` is true")
    }

    return nil
}
```

**ForceNew Logic Issues:**
```go
// Debug ForceNew conditions
pluginsdk.ForceNewIfChange("field_name", func(ctx context.Context, old, new, meta interface{}) bool {
    log.Printf("[DEBUG] ForceNew check: old=%v, new=%v", old, new)

    shouldForceNew := old.(string) != new.(string)
    log.Printf("[DEBUG] ForceNew result: %t", shouldForceNew)

    return shouldForceNew
}),
```

### Azure API Integration Issues

**Client Configuration Problems:**
```go
// Debug client initialization
func debugClientSetup(metadata sdk.ResourceMetaData) {
    log.Printf("[DEBUG] Subscription ID: %s", metadata.Client.Account.SubscriptionId)
    log.Printf("[DEBUG] Client features: %+v", metadata.Client.Features)

    // Test client connectivity
    client := metadata.Client.ServiceName.ResourceClient
    // Make a lightweight API call to test
}
```

**Resource ID Parsing Issues:**
```go
// Debug resource ID parsing
id, err := parse.ServiceNameID(resourceId)
if err != nil {
    log.Printf("[DEBUG] Failed to parse resource ID '%s': %+v", resourceId, err)
    return fmt.Errorf("parsing Resource ID `%s`: %+v", resourceId, err)
}
log.Printf("[DEBUG] Parsed ID: %+v", id)
```
---
[⬆️ Back to top](#🔧-troubleshooting-decision-trees)
