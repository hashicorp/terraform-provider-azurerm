---
applyTo: "internal/**/*.go"
description: Error handling patterns and standards for the Terraform AzureRM provider including message formatting, error types, and debugging guidelines.
---

# Error Handling Patterns

Quick navigation: [💬 Error Messages](#-error-message-standards) | [🔍 Error Types](#-error-type-patterns) | [🐛 Debugging](#-debugging-patterns) | [🔄 State Management](#-state-management-errors)

## 💬 Error Message Standards

### Field Names and Values with Backticks

**Field names and values must be wrapped in backticks for clarity:**

```go
// GOOD - Field names and values properly formatted with backticks
return fmt.Errorf("creating Storage Account `%s` with SKU `%s` in location `%s`: %+v", name, skuName, location, err)
return fmt.Errorf("property `account_tier` must be `Standard` or `Premium`, got `%s`", accountTier)
return fmt.Errorf("field `zones` cannot be set when `availability_set_id` is specified")

// BAD - Missing backticks around field names and values
return fmt.Errorf("creating Storage Account %q with SKU %s in location %s: %+v", name, skuName, location, err)
return fmt.Errorf("property account_tier must be Standard or Premium, got %s", accountTier)
return fmt.Errorf("field zones cannot be set when availability_set_id is specified")
```

### Lowercase, No Punctuation, Descriptive

**Error messages must follow Go standards:**

```go
// GOOD - Lowercase, no punctuation, descriptive error messages
return fmt.Errorf("creating resource group `%s` in location `%s`: %+v", name, location, err)
return fmt.Errorf("updating virtual network `%s`: %+v", id, err)

// BAD - Incorrect casing, punctuation, or vague messages
return fmt.Errorf("Creating Resource Group %q in Location %q: %v", name, location, err)
return fmt.Errorf("error updating virtual network: %s", err.Error())
```

### Verbose Error Formatting

**Always use `%+v` for comprehensive error context:**

```go
// GOOD - Verbose error formatting provides full context
return fmt.Errorf("creating CDN Front Door Profile `%s`: %+v", name, err)
return fmt.Errorf("updating Network Security Group rules: %+v", err)
return fmt.Errorf("polling for completion of operation: %+v", err)

// BAD - Simple formatting loses error context
return fmt.Errorf("creating CDN Front Door Profile `%s`: %v", name, err)
return fmt.Errorf("updating Network Security Group rules: %s", err.Error())
return fmt.Errorf("polling for completion of operation: %w", err)
```

### Clear Context and Actionable Information

```go
// GOOD - Clear context and actionable information
return fmt.Errorf("creating Storage Account `%s`: name must be globally unique, try a different name: %+v", name, err)
return fmt.Errorf("VM size `%s` is not available in location `%s`, choose a different size or location", size, location)
return fmt.Errorf("property `disk_size_gb` must be between 1 and 32767, got %d", diskSize)

// BAD - Vague, unhelpful error messages
return fmt.Errorf("creating Storage Account failed: %+v", err)
return fmt.Errorf("VM size problem: %+v", err)
return fmt.Errorf("invalid disk size: %+v", err)
```

### Contractions Policy

**Do not use contractions in error messages:**

```go
// GOOD - Full words
return fmt.Errorf("property `name` cannot be empty")
return fmt.Errorf("resource `%s` is not available in this region", resourceName)
return fmt.Errorf("field `enabled` cannot be disabled once set to true")

// BAD - Contractions
return fmt.Errorf("property `name` can't be empty")
return fmt.Errorf("resource `%s` isn't available in this region", resourceName)
return fmt.Errorf("field `enabled` can't be disabled once set to true")
```

## 🔍 Error Type Patterns

### Typed Resource Error Patterns

```go
// Use metadata.Decode for model decoding errors
var model ServiceNameModel
if err := metadata.Decode(&model); err != nil {
    return fmt.Errorf("decoding: %+v", err)
}

// Use metadata.Logger for structured logging
metadata.Logger.Infof("Import check for %s", id)

// Use metadata.ResourceRequiresImport for import conflicts
if !response.WasNotFound(existing.HttpResponse) {
    return metadata.ResourceRequiresImport(r.ResourceType(), id)
}

// Use metadata.MarkAsGone for deleted resources
if response.WasNotFound(resp.HttpResponse) {
    return metadata.MarkAsGone(id)
}

// Use metadata.SetID for resource ID management
metadata.SetID(id)

// Use metadata.Encode for state management
return metadata.Encode(&model)
```

### UnTyped Resource Error Patterns

```go
// Use consistent error formatting with context
if err != nil {
    return fmt.Errorf("creating Resource `%s`: %+v", name, err)
}

// Include resource information in error messages
if response.WasNotFound(resp.HttpResponse) {
    log.Printf("[DEBUG] Resource `%s` was not found - removing from state", id.ResourceName)
    d.SetId("")
    return nil
}

// Handle Azure-specific errors
if response.WasThrottled(resp.HttpResponse) {
    return resource.RetryableError(fmt.Errorf("request was throttled"))
}
```

### Resource Not Found Messaging

```go
// Typed resource approach
if response.WasNotFound(resp.HttpResponse) {
    return metadata.MarkAsGone(id)
}

// UnTyped resource approach
if response.WasNotFound(resp.HttpResponse) {
    log.Printf("[DEBUG] Storage Account `%s` was not found - removing from state", id.StorageAccountName)
    d.SetId("")
    return nil
}

// Data source approach (should return error, not mark as gone)
if response.WasNotFound(resp.HttpResponse) {
    return fmt.Errorf("CDN Front Door Profile `%s` was not found in Resource Group `%s`", profileName, resourceGroupName)
}
```

### Parsing Error Context

```go
// GOOD - Clear parsing error context
id, err := parse.VirtualMachineID(d.Id())
if err != nil {
    return fmt.Errorf("parsing Virtual Machine ID `%s`: %+v", d.Id(), err)
}

// Typed resource approach
id, err := parse.ServiceNameID(metadata.ResourceData.Id())
if err != nil {
    return fmt.Errorf("parsing Service ID `%s`: %+v", metadata.ResourceData.Id(), err)
}
```

## 🐛 Debugging Patterns

### PATCH Operation + "None" Pattern Debugging

**Common Symptoms:**
- Resource state shows fields as disabled, but Azure portal shows them as enabled
- Tests pass on creation but fail when testing disable → re-enable scenarios
- Azure API calls return success, but resource configuration doesn't change
- Residual state persists after removing Terraform configuration blocks

**Root Cause Analysis Framework:**

1. **Identify the HTTP Method**: Check if the Azure service uses PATCH vs PUT operations
   ```powershell
   # Look for PatchThenPoll vs CreateOrUpdateThenPoll in Azure SDK calls
   grep -r "PatchThenPoll|CreateOrUpdateThenPoll" internal/services/servicename/
   ```

2. **Trace Azure SDK Filtering**: Verify if nil values are being filtered out
   ```go
   // Look for patterns like this that cause issues:
   if len(input) == 0 {
       return nil // SDK filters this out, Azure never gets disable command
   }
   ```

3. **Check "None" Pattern Implementation**: Ensure disabled features are explicit
   ```go
   // WRONG - Causes residual state
   func ExpandFeature(input []interface{}) *azuretype.Feature {
       if len(input) == 0 {
           return nil
       }
       // Configure only enabled features
   }

   // RIGHT - Prevents residual state
   func ExpandFeature(input []interface{}) *azuretype.Feature {
       result := &azuretype.Feature{
           Enabled: pointer.To(false), // Explicit disable
       }
       if len(input) > 0 {
           result.Enabled = pointer.To(true)
       }
       return result
   }
   ```

### Console Line Wrapping Detection

**🚨 CRITICAL: Console Line Wrapping Detection Protocol 🚨**

**CONSOLE LINE WRAPPING WARNING**: When reviewing git diff output in terminal/console, be aware that long lines may wrap and appear malformed. Always verify actual file content for syntax validation, especially for JSON, YAML, or structured data files.

**VERIFICATION PROTOCOL FOR SUSPECTED ISSUES**:

🔍 **MANDATORY VERIFICATION STEPS:**
1. **STOP** - If text appears broken/fragmented, this is likely console wrapping
2. **VERIFY** - Use `Get-Content filename` (PowerShell) or `cat filename` (bash) to check actual file content
3. **VALIDATE** - For JSON/structured files: `Get-Content file.json | ConvertFrom-Json` or `jq "." file.json`

### 🚨 **Console Wrapping Red Flags:**
- ❌ Text breaks mid-sentence or mid-word without logical reason
- ❌ Missing closing quotes/brackets that don't make sense contextually
- ❌ Fragmented lines that appear to continue elsewhere in the diff
- ❌ Content looks syntactically invalid but conceptually correct
- ❌ Long lines in git diff output that suddenly break

#### ✅ **GOLDEN RULE**: If actual file content is valid → acknowledge console wrapping → do NOT flag as corruption

## 🔄 State Management Errors

### Import Conflict Detection

```go
// Typed resource import conflict
if !response.WasNotFound(existing.HttpResponse) {
    return metadata.ResourceRequiresImport(r.ResourceType(), id)
}

// UnTyped resource import conflict
if existing.StatusCode != http.StatusNotFound {
    return tf.ImportAsExistsError("azurerm_resource", id.ID())
}
```

### State Validation Errors

```go
// Ensure required model fields are populated
if model == nil {
    return fmt.Errorf("retrieving %s: model was nil", id)
}

if props := model.Properties; props == nil {
    return fmt.Errorf("retrieving %s: properties was nil", id)
}
```

### Timeout Error Handling

```go
// Use context-aware timeout errors
ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
defer cancel()

if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
    select {
    case <-ctx.Done():
        return fmt.Errorf("creating %s: operation timed out after 30 minutes", id)
    default:
        return fmt.Errorf("creating %s: %+v", id, err)
    }
}
```

## 🚨 Common Error Scenarios

### Azure API Rate Limiting

```go
// Exponential backoff for throttled requests
if response.WasThrottled(resp.HttpResponse) {
    return resource.RetryableError(fmt.Errorf("request was throttled, retrying"))
}

// Check for specific throttling error codes
if strings.Contains(err.Error(), "TooManyRequests") {
    return resource.RetryableError(fmt.Errorf("Azure API rate limit exceeded, retrying: %+v", err))
}
```

### Azure Resource Dependencies

```go
// Handle dependency conflicts
if strings.Contains(err.Error(), "ResourceGroupBeingDeleted") {
    return resource.RetryableError(fmt.Errorf("resource group is being deleted, retrying: %+v", err))
}

// Handle resource locks
if strings.Contains(err.Error(), "ScopeLocked") {
    return fmt.Errorf("resource `%s` is locked and cannot be modified: %+v", id, err)
}
```

### Azure Service Quotas

```go
// Handle quota exceeded errors
if strings.Contains(err.Error(), "QuotaExceeded") {
    return fmt.Errorf("Azure service quota exceeded for resource `%s`: increase quota or use a different region: %+v", id, err)
}

// Handle specific quota types
if strings.Contains(err.Error(), "cores quota") {
    return fmt.Errorf("CPU cores quota exceeded in region `%s`: request quota increase in Azure portal: %+v", location, err)
}
```

### Validation Error Examples

```go
// CustomizeDiff validation errors
func validateConfiguration(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
    if diff.Get("enabled").(bool) && diff.Get("configuration") == nil {
        return fmt.Errorf("`configuration` is required when `enabled` is true")
    }

    if diff.Get("sku_name").(string) == "Premium" && !diff.Get("zone_redundant").(bool) {
        return fmt.Errorf("`zone_redundant` must be true when `sku_name` is `Premium`")
    }

    return nil
}

// Schema validation errors
func ValidateResourceName(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)

    if len(value) > 64 {
        errors = append(errors, fmt.Errorf("property `%s` cannot be longer than 64 characters, got %d", k, len(value)))
    }

    if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(value) {
        errors = append(errors, fmt.Errorf("property `%s` can only contain alphanumeric characters and hyphens", k))
    }

    return warnings, errors
}
```

## 🏗️ Error Recovery Patterns

### Graceful Degradation

```go
// Handle optional features gracefully
func expandOptionalFeature(input []interface{}) *azureapi.OptionalFeature {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[WARN] Failed to expand optional feature, using defaults: %v", r)
        }
    }()

    if len(input) == 0 {
        return nil
    }

    // Process optional feature
    return processFeature(input)
}
```

### Retry Logic with Exponential Backoff

```go
func retryWithExponentialBackoff(ctx context.Context, operation func() error, logger interface{}) error {
    const maxRetries = 5
    const baseDelay = 1 * time.Second
    const maxDelay = 32 * time.Second

    for attempt := 0; attempt < maxRetries; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }

        // Check if it's a retryable error
        if !isRetryableError(err) {
            return err
        }

        if attempt == maxRetries-1 {
            return fmt.Errorf("operation failed after %d attempts: %+v", maxRetries, err)
        }

        // Calculate exponential backoff delay
        delay := time.Duration(math.Pow(2, float64(attempt))) * baseDelay
        if delay > maxDelay {
            delay = maxDelay
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(delay):
            continue
        }
    }

    return fmt.Errorf("operation failed after %d attempts", maxRetries)
}

func isRetryableError(err error) bool {
    if err == nil {
        return false
    }

    errStr := strings.ToLower(err.Error())
    retryableErrors := []string{
        "throttled",
        "toomanyrequests",
        "internalservererror",
        "serviceunavailable",
        "timeout",
    }

    for _, retryableErr := range retryableErrors {
        if strings.Contains(errStr, retryableErr) {
            return true
        }
    }

    return false
}
```

---

## Quick Reference Links

- 🏗️ **Main Implementation Guide**: [implementation-guide.md](./implementation-guide.md)
- ⚡ **Azure Patterns**: [azure-patterns.md](./azure-patterns.md)
- 🧪 **Testing Guide**: [testing-guidelines.md](./testing-guidelines.md)
- 📝 **Documentation Guide**: [documentation-guidelines.md](./documentation-guidelines.md)
