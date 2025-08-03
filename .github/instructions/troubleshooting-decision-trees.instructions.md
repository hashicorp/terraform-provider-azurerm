---
applyTo: "internal/**/*.go"
description: Troubleshooting decision trees and diagnostic patterns for the Terraform AzureRM provider including common issues, root cause analysis, and resolution frameworks.
---

# Troubleshooting Decision Trees

Troubleshooting decision trees and diagnostic patterns for the Terraform AzureRM provider including common issues, root cause analysis, and resolution frameworks.

**Quick navigation:** [🔧 Common Issues](#🔧-common-issue-resolution-flowchart) | [🔍 Root Cause Analysis](#🔍-root-cause-analysis-framework) | [🚨 Error Diagnostics](#🚨-error-diagnostic-patterns) | [🔄 State Issues](#🔄-state-management-troubleshooting) | [🏗️ Implementation Choice](#🔧-implementation-choice-decision-trees) | [🧪 Testing Strategy](#🧪-testing-strategy-decision-trees) | [📝 Code Quality](#📝-code-quality-decision-trees) | [🔄 Azure Lifecycle](#🔄-azure-resource-lifecycle-decision-trees) | [🔍 Performance](#🔍-performance-optimization-decision-trees)

## 🔧 Common Issue Resolution Flowchart

### Resource Creation Failures

```text
Resource Creation Fails
├─ Azure API Error 409 (Conflict)
│  ├─ Check import conflict detection
│  ├─ Verify resource doesn't already exist
│  └─ Review RequiresImport implementation
├─ Azure API Error 400 (Bad Request)
│  ├─ Validate field combinations in CustomizeDiff
│  ├─ Check required field validation
│  └─ Verify Azure SDK parameter mapping
├─ Timeout Error
│  ├─ Increase timeout values for long-running operations
│  ├─ Check Azure service health
│  └─ Verify polling implementation for LROs
├─ Permission Error (403)
│  ├─ Verify service principal permissions
│  ├─ Check Azure RBAC assignments
│  └─ Validate subscription access
└─ Unknown Error
   ├─ Check Azure SDK version compatibility
   ├─ Review Azure service API changes
   └─ Validate authentication configuration
```

### PATCH Operation Issues

```text
PATCH Operation Problems
├─ Residual State (features remain enabled after removal)
│  ├─ Check "None" pattern implementation
│  ├─ Verify explicit disable commands
│  └─ Review Azure SDK nil filtering behavior
├─ Fields Not Updating
│  ├─ Verify expand function completeness
│  ├─ Check Azure API field mapping
│  └─ Validate pointer usage patterns
├─ State Drift Detection
│  ├─ Check flatten function accuracy
│  ├─ Verify Read function implementation
│  └─ Review computed field handling
└─ Import Failures
   ├─ Verify resource ID parsing
   ├─ Check flatten function completeness
   └─ Validate state reconstruction logic
```

### Authentication and Authorization Issues

```text
Authentication Problems
├─ Invalid Credentials
│  ├─ Verify environment variables are set
│  ├─ Check credential format validation
│  └─ Test authentication outside Terraform
├─ Token Expiration
│  ├─ Implement token refresh logic
│  ├─ Check token lifetime settings
│  └─ Verify refresh token handling
├─ Insufficient Permissions
│  ├─ Review required Azure permissions
│  ├─ Check resource group access
│  └─ Validate subscription-level permissions
└─ Multi-Tenant Issues
   ├─ Verify tenant ID configuration
   ├─ Check cross-tenant access
   └─ Review guest user permissions
```
---
[⬆️ Back to top](#troubleshooting-decision-trees)

## 🔍 Root Cause Analysis Framework

### Systematic Debugging Approach

```go
func debugResourceIssue(ctx context.Context, resourceType string, operation string) {
    logger := log.WithFields(logrus.Fields{
        "resource_type": resourceType,
        "operation":     operation,
        "debug_session": generateDebugID(),
    })

    // Step 1: Environment validation
    if err := validateEnvironment(); err != nil {
        logger.Errorf("Environment validation failed: %+v", err)
        return
    }

    // Step 2: Configuration analysis
    if err := analyzeConfiguration(); err != nil {
        logger.Errorf("Configuration analysis failed: %+v", err)
        return
    }

    // Step 3: Azure API testing
    if err := testAzureAPI(ctx); err != nil {
        logger.Errorf("Azure API test failed: %+v", err)
        return
    }

    // Step 4: State comparison
    if err := compareExpectedVsActualState(); err != nil {
        logger.Errorf("State comparison failed: %+v", err)
        return
    }

    logger.Info("Debug analysis complete")
}
```

### Configuration Validation

```go
func analyzeConfiguration() error {
    checks := []struct {
        name string
        fn   func() error
    }{
        {"Schema Validation", validateSchemaConfiguration},
        {"Field Dependencies", validateFieldDependencies},
        {"Azure Constraints", validateAzureConstraints},
        {"Resource Limits", validateResourceLimits},
    }

    for _, check := range checks {
        if err := check.fn(); err != nil {
            return fmt.Errorf("%s failed: %+v", check.name, err)
        }
    }

    return nil
}

func validateFieldDependencies() error {
    // Check for missing required field combinations
    // Validate conditional field requirements
    // Verify CustomizeDiff logic alignment
    return nil
}
```

### Azure API Diagnostics

```go
func testAzureAPI(ctx context.Context) error {
    // Test basic connectivity
    if err := testConnectivity(ctx); err != nil {
        return fmt.Errorf("connectivity test failed: %+v", err)
    }

    // Test authentication
    if err := testAuthentication(ctx); err != nil {
        return fmt.Errorf("authentication test failed: %+v", err)
    }

    // Test specific API endpoints
    if err := testAPIEndpoints(ctx); err != nil {
        return fmt.Errorf("API endpoint test failed: %+v", err)
    }

    return nil
}

func testAPIEndpoints(ctx context.Context) error {
    endpoints := []string{
        "/subscriptions/{subscriptionId}/providers/Microsoft.Resources/resourceGroups",
        "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProvider}",
    }

    for _, endpoint := range endpoints {
        if err := testEndpoint(ctx, endpoint); err != nil {
            return fmt.Errorf("endpoint %s failed: %+v", endpoint, err)
        }
    }

    return nil
}
```
---
[⬆️ Back to top](#troubleshooting-decision-trees)

## 🚨 Error Diagnostic Patterns

### Error Classification System

```go
type ErrorCategory int

const (
    AuthenticationError ErrorCategory = iota
    AuthorizationError
    ConfigurationError
    AzureAPIError
    NetworkError
    TimeoutError
    StateError
    UnknownError
)

func classifyError(err error) ErrorCategory {
    if err == nil {
        return UnknownError
    }

    errorString := strings.ToLower(err.Error())

    switch {
    case strings.Contains(errorString, "unauthorized") || strings.Contains(errorString, "authentication"):
        return AuthenticationError
    case strings.Contains(errorString, "forbidden") || strings.Contains(errorString, "permission"):
        return AuthorizationError
    case strings.Contains(errorString, "bad request") || strings.Contains(errorString, "invalid"):
        return ConfigurationError
    case strings.Contains(errorString, "timeout") || strings.Contains(errorString, "deadline"):
        return TimeoutError
    case strings.Contains(errorString, "conflict") || strings.Contains(errorString, "already exists"):
        return StateError
    default:
        return AzureAPIError
    }
}
```

### Error Resolution Mapping

```go
func getResolutionSteps(category ErrorCategory, err error) []string {
    resolutions := map[ErrorCategory][]string{
        AuthenticationError: {
            "Verify ARM_CLIENT_ID, ARM_CLIENT_SECRET, ARM_TENANT_ID environment variables",
            "Check service principal credentials",
            "Test authentication with Azure CLI: az login --service-principal",
            "Verify credential expiration dates",
        },
        AuthorizationError: {
            "Check Azure RBAC role assignments",
            "Verify service principal has required permissions",
            "Review resource group access permissions",
            "Check subscription-level permissions",
        },
        ConfigurationError: {
            "Run CustomizeDiff validation tests",
            "Verify field combinations and dependencies",
            "Check Azure service constraints",
            "Validate input data types and formats",
        },
        TimeoutError: {
            "Increase timeout values in resource configuration",
            "Check Azure service health status",
            "Verify network connectivity to Azure endpoints",
            "Review long-running operation polling implementation",
        },
        StateError: {
            "Check for existing resources with same name",
            "Verify import detection logic",
            "Review state file for conflicts",
            "Validate resource ID uniqueness",
        },
    }

    return resolutions[category]
}
```

### Automated Diagnostics

```go
func runAutomatedDiagnostics(ctx context.Context, resource interface{}) DiagnosticReport {
    report := DiagnosticReport{
        Timestamp: time.Now(),
        Checks:    make(map[string]CheckResult),
    }

    // Configuration checks
    report.Checks["schema_validation"] = validateResourceSchema(resource)
    report.Checks["field_dependencies"] = validateFieldDependencies(resource)

    // Azure connectivity checks
    report.Checks["azure_connectivity"] = testAzureConnectivity(ctx)
    report.Checks["authentication"] = testAuthentication(ctx)

    // Performance checks
    report.Checks["api_latency"] = measureAPILatency(ctx)
    report.Checks["resource_quotas"] = checkResourceQuotas(ctx)

    // State checks
    report.Checks["state_consistency"] = validateStateConsistency(resource)

    return report
}

type DiagnosticReport struct {
    Timestamp time.Time                `json:"timestamp"`
    Checks    map[string]CheckResult   `json:"checks"`
    Summary   string                   `json:"summary"`
}

type CheckResult struct {
    Status  string      `json:"status"` // "pass", "fail", "warning"
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}
```
---
[⬆️ Back to top](#troubleshooting-decision-trees)

## 🔄 State Management Troubleshooting

### State Drift Detection

```text
State Drift Issues
├─ Read Function Problems
│  ├─ Check API response parsing
│  ├─ Verify flatten function accuracy
│  └─ Review null/empty value handling
├─ Azure Resource Changes
│  ├─ Check for manual Azure portal changes
│  ├─ Verify Azure policy effects
│  └─ Review Azure automation impacts
├─ Provider Version Changes
│  ├─ Check for breaking changes in provider updates
│  ├─ Review schema modifications
│  └─ Validate migration requirements
└─ Terraform State Corruption
   ├─ Backup and restore state file
   ├─ Use terraform state pull/push commands
   └─ Consider terraform refresh operations
```

### Import Issues

```go
func debugImportIssues(resourceID string) error {
    // Step 1: Validate resource ID format
    if err := validateResourceIDFormat(resourceID); err != nil {
        return fmt.Errorf("invalid resource ID format: %+v", err)
    }

    // Step 2: Check resource existence in Azure
    exists, err := checkResourceExistence(resourceID)
    if err != nil {
        return fmt.Errorf("error checking resource existence: %+v", err)
    }
    if !exists {
        return fmt.Errorf("resource does not exist in Azure")
    }

    // Step 3: Test resource parsing
    if err := testResourceParsing(resourceID); err != nil {
        return fmt.Errorf("resource parsing failed: %+v", err)
    }

    // Step 4: Validate flatten functions
    if err := testFlattenFunctions(resourceID); err != nil {
        return fmt.Errorf("flatten function validation failed: %+v", err)
    }

    return nil
}
```

### State Reconstruction

```go
func reconstructResourceState(ctx context.Context, resourceID string) (map[string]interface{}, error) {
    // Parse resource ID
    id, err := parseResourceID(resourceID)
    if err != nil {
        return nil, fmt.Errorf("parsing resource ID: %+v", err)
    }

    // Fetch current state from Azure
    azureState, err := fetchFromAzure(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("fetching from Azure: %+v", err)
    }

    // Apply flatten functions
    terraformState := flattenToTerraformState(azureState)

    // Validate reconstructed state
    if err := validateReconstructedState(terraformState); err != nil {
        return nil, fmt.Errorf("state validation failed: %+v", err)
    }

    return terraformState, nil
}
```
---
[⬆️ Back to top](#troubleshooting-decision-trees)

## 🔧 Implementation Choice Decision Trees

### Implementation Approach Selection

```text
Need to implement new resource/data source?
├─ NEW implementation
│  ├─ Use Typed Resource Implementation (Preferred)
│  ├─ Benefits: Type safety, better error handling, metadata
│  └─ Pattern: sdk.Resource with receiver methods
├─ EXISTING resource maintenance
│  ├─ Continue with Untyped Resource Implementation
│  ├─ Maintain existing function-based CRUD patterns
│  └─ Pattern: pluginsdk.Resource with function pointers
├─ Major refactor/migration
│  ├─ Consider migration to Typed Implementation
│  ├─ Evaluate cost/benefit of migration
│  └─ Follow migration guide patterns
└─ Bug fix/minor change
   ├─ Maintain existing implementation approach
   ├─ Don't mix typed/untyped patterns
   └─ Apply same standards regardless of approach
```

### Schema Design Decision Tree

```text
Designing resource schema?
├─ Azure API has wrapper structures
│  ├─ Single-purpose wrapper? → Consider schema flattening
│  ├─ Logical grouping? → Maintain nested structure
│  ├─ User experience improved? → Flatten responsibly
│  └─ Complex validation needed? → Keep structure for clarity
├─ Field validation requirements
│  ├─ Simple validation? → Use schema ValidateFunc
│  ├─ Azure service constraint? → Use CustomizeDiff
│  ├─ Field combinations? → Use CustomizeDiff with GetRawConfig
│  └─ Complex state transitions? → Use CustomizeDiff ForceNew
├─ Azure PATCH operations
│  ├─ Service uses PATCH? → Implement "None" pattern
│  ├─ Residual state concerns? → Explicit disable patterns
│  ├─ Feature toggles? → Always return complete structures
│  └─ nil filtering issues? → Use pointer.To(false) for disable
└─ Optional field handling
   ├─ Go zero value conflicts? → Use GetRawConfig().IsNull()
   ├─ Required combinations? → Validate in CustomizeDiff
   ├─ Conditional requirements? → Check field existence first
   └─ Performance critical? → Minimize raw config access
```
---
[⬆️ Back to top](#troubleshooting-decision-trees)

## 🧪 Testing Strategy Decision Trees

### Test Type Selection

```text
What needs testing?
├─ Unit testing required
│  ├─ Parser functions → Table-driven tests
│  ├─ Validation logic → Error/success scenarios
│  ├─ Utility functions → Edge cases and boundaries
│  └─ SDK integration → Mock Azure responses
├─ Acceptance testing required
│  ├─ Basic CRUD → ExistsInAzure + ImportStep pattern
│  ├─ Update scenarios → Multi-step resource lifecycle
│  ├─ Error scenarios → ExpectError with regexp
│  ├─ CustomizeDiff validation → MANDATORY error testing
│  └─ Import functionality → RequiresImport pattern
├─ Data source testing
│  ├─ Field validation → Key checks (NOT redundant with ImportStep)
│  ├─ Computed attributes → Verify population
│  ├─ Complex structures → Nested value validation
│  └─ No ImportStep → All validation must be explicit
└─ Performance testing
   ├─ Large resource sets → Parallel processing patterns
   ├─ Azure API limits → Rate limiting and retry logic
   ├─ Long-running operations → Timeout and polling
   └─ Memory usage → Large state file handling
```

### Test Execution Decision Tree

```text
Ready to run tests?
├─ ⚠️ STOP: Never run automatically
│  ├─ Tests create REAL Azure resources
│  ├─ Require valid Azure credentials
│  ├─ Generate billable charges
│  └─ Need user confirmation
├─ Manual execution only
│  ├─ Provide exact command to user
│  ├─ Explain purpose and duration
│  ├─ List prerequisites (credentials)
│  └─ Warn about Azure costs
├─ Environment check
│  ├─ ARM_SUBSCRIPTION_ID set?
│  ├─ ARM_CLIENT_ID/SECRET/TENANT_ID configured?
│  ├─ Azure permissions sufficient?
│  └─ Test region availability confirmed?
└─ Cleanup verification
   ├─ Provider features for force deletion
   ├─ Resource dependencies handled
   ├─ Soft-delete considerations
   └─ Billing impact minimized
```
---
[⬆️ Back to top](#troubleshooting-decision-trees)

## 📝 Code Quality Decision Trees

### Comment Policy Decision Tree

```text
About to add a comment?
├─ ⚠️ MANDATORY STOP: Zero tolerance policy
│  ├─ Can code be self-explanatory instead?
│  ├─ Better naming eliminate need?
│  ├─ Function extraction possible?
│  └─ Structure reorganization help?
├─ Exception evaluation (4 cases only)
│  ├─ Azure API quirk not obvious? → MAY be acceptable
│  ├─ Complex business logic that can't be simplified? → MAY be acceptable
│  ├─ Azure SDK workaround/limitation? → MAY be acceptable
│  ├─ Non-obvious state pattern (PATCH, residual state)? → MAY be acceptable
│  └─ Everything else → NO COMMENT (refactor instead)
├─ Justification required
│  ├─ Which exception case applies?
│  ├─ Why can't code be self-explanatory?
│  ├─ What specific Azure behavior needs documentation?
│  └─ Can this comment be eliminated through better code?
└─ Final check
   ├─ Is this truly necessary?
   ├─ Does it add value beyond code?
   ├─ Will future developers need this context?
   └─ Can refactoring eliminate this need?
```

### Error Handling Decision Tree

```text
Handling Azure API errors?
├─ Resource not found (404)
│  ├─ During Read operation? → metadata.MarkAsGone() or return nil
│  ├─ During Create operation? → Proceed with creation
│  ├─ During Update operation? → Resource deleted externally
│  └─ During Delete operation? → Consider already deleted
├─ Authentication errors (401)
│  ├─ Check credential configuration
│  ├─ Verify service principal permissions
│  ├─ Test token expiration
│  └─ Validate tenant access
├─ Authorization errors (403)
│  ├─ Check RBAC role assignments
│  ├─ Verify resource group permissions
│  ├─ Check subscription access
│  └─ Review resource-specific permissions
├─ Conflict errors (409)
│  ├─ Import conflict? → Return metadata.ResourceRequiresImport
│  ├─ Resource state conflict? → Check for external changes
│  ├─ Azure policy violation? → Review compliance requirements
│  └─ Concurrent modification? → Implement retry logic
├─ Rate limiting (429)
│  ├─ Implement exponential backoff
│  ├─ Check for batch operations
│  ├─ Review concurrent request patterns
│  └─ Consider Azure SDK automatic retry
└─ Unknown errors
   ├─ Log full error context
   ├─ Include operation details
   ├─ Preserve Azure request ID
   └─ Provide actionable error message
```
---
[⬆️ Back to top](#troubleshooting-decision-trees)

## 🔄 Azure Resource Lifecycle Decision Trees

### PATCH Operation Troubleshooting

```text
Azure resource not updating correctly?
├─ Service uses PATCH operations?
│  ├─ Features remain enabled after removal? → Implement explicit disable
│  ├─ nil values being filtered? → Use pointer.To(false) for disabled features
│  ├─ Residual state persisting? → Return complete structure with all features
│  └─ Wrapper structures causing issues? → Apply "None" pattern
├─ Configuration removal not working?
│  ├─ Check expand function returns disabled state
│  ├─ Verify all features explicitly set to false
│  ├─ Ensure required fields included even when disabled
│  └─ Test with empty configuration scenarios
├─ State drift detection failing?
│  ├─ Flatten function handling disabled features?
│  ├─ Read operation detecting all state changes?
│  ├─ Computed fields properly updated?
│  └─ Import functionality reconstructing full state?
└─ Field combination validation?
   ├─ CustomizeDiff handling PATCH requirements?
   ├─ GetRawConfig() used for field existence checks?
   ├─ Zero value validation preventing false errors?
   └─ Conditional logic matching Azure API constraints?
```

### CustomizeDiff Troubleshooting

```text
CustomizeDiff validation not working?
├─ Import requirements issues?
│  ├─ Typed resource? → May need dual imports (schema + pluginsdk)
│  ├─ Untyped resource? → Usually only pluginsdk import sufficient
│  ├─ Function signature mismatch? → Check *schema.ResourceDiff vs *pluginsdk.ResourceDiff
│  └─ Compilation errors? → Verify import requirements
├─ Zero value validation problems?
│  ├─ Optional fields causing false errors? → Use GetRawConfig().IsNull()
│  ├─ Go zero values triggering validation? → Check field existence first
│  ├─ Required fields incorrectly validated? → Use diff.Get() for required fields
│  └─ Performance issues? → Minimize raw config access overhead
├─ Field removal ForceNew not working?
│  ├─ SetNew() called before ForceNew()? → Both required for visibility
│  ├─ Plan showing state change? → SetNew creates visible transition
│  ├─ Update-only logic? → Check diff.Id() != "" to avoid creation issues
│  └─ Error handling? → Wrap SetNew errors with descriptive context
├─ Boolean expressions causing linting errors?
│  ├─ Using verbose comparisons? → Simplify: old.(bool) && !new.(bool)
│  ├─ Explicit true/false checks? → Use direct boolean semantics
│  ├─ Compliance with gosimple? → Apply simplified expressions
│  └─ Readability maintained? → Shorter expressions are clearer
└─ Testing validation logic?
   ├─ Error scenarios covered? → Use ExpectError with regexp.MustCompile()
   ├─ Success scenarios tested? → Basic/update/complete tests handle these
   ├─ Edge cases included? → Test boundary conditions
   └─ Azure API constraints validated? → Match actual service behavior
```
---
[⬆️ Back to top](#troubleshooting-decision-trees)

## 🔍 Performance Optimization Decision Trees

### Resource Management Performance

```text
Performance issues with resource operations?
├─ Create/Update operations slow?
│  ├─ Long-running operations? → Verify polling implementation
│  ├─ Multiple API calls? → Consider batch operations
│  ├─ Timeout values appropriate? → Match Azure service SLA
│  └─ Connection pooling configured? → Check HTTP client settings
├─ Read operations inefficient?
│  ├─ Unnecessary API calls? → Cache computed values appropriately
│  ├─ Large response payloads? → Filter required fields only
│  ├─ State comparison overhead? → Optimize flatten functions
│  └─ Drift detection expensive? → Implement incremental checks
├─ Multiple resources slow?
│  ├─ Sequential processing? → Implement parallel patterns with semaphores
│  ├─ Resource dependencies? → Optimize dependency resolution
│  ├─ Azure API rate limits? → Implement intelligent backoff
│  └─ Memory usage high? → Optimize state management patterns
├─ Import operations problematic?
│  ├─ Resource ID parsing expensive? → Cache parsed components
│  ├─ State reconstruction slow? → Optimize flatten functions
│  ├─ Multiple resources import? → Batch processing where possible
│  └─ Error handling overhead? → Streamline validation logic
└─ Test execution performance?
   ├─ Acceptance tests slow? → Minimize Azure resource creation
   ├─ Setup/teardown expensive? → Optimize test fixtures
   ├─ Parallel test conflicts? → Review resource naming/isolation
   └─ CI/CD pipeline slow? → Consider test parallelization
```
---
[⬆️ Back to top](#troubleshooting-decision-trees)

## Quick Reference Links

- 🏠 **Home**: [../copilot-instructions.md](../copilot-instructions.md)
- ☁️ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md)
- 🏗️ **Implementation Guide**: [implementation-guide.instructions.md](./implementation-guide.instructions.md)
- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md)
- ❌ **Error Patterns**: [error-patterns.instructions.md](./error-patterns.instructions.md)
- 🔄 **Migration Guide**: [migration-guide.instructions.md](./migration-guide.instructions.md)
- 🏢 **Provider Guidelines**: [provider-guidelines.instructions.md](./provider-guidelines.instructions.md)
- 📐 **Schema Patterns**: [schema-patterns.instructions.md](./schema-patterns.instructions.md)
- 📋 **Code Clarity**: [code-clarity-enforcement.instructions.md](./code-clarity-enforcement.instructions.md)

### 🚀 Enhanced Guidance Files

- 🔄 **API Evolution**: [api-evolution-patterns.instructions.md](./api-evolution-patterns.instructions.md)
- ⚡ **Performance**: [performance-optimization.instructions.md](./performance-optimization.instructions.md)
- 🔐 **Security**: [security-compliance.instructions.md](./security-compliance.instructions.md)

---
[⬆️ Back to top](#troubleshooting-decision-trees)
