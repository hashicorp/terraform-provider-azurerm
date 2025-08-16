---
applyTo: "internal/**/*.go"
description: "This is the official Terraform Provider for Azure (Resource Manager), written in Go. It enables Terraform to manage Azure resources through the Azure Resource Manager APIs."
---
# Custom instructions

This is the official Terraform Provider for Azure (Resource Manager), written in Go. It enables Terraform to manage Azure resources through the Azure Resource Manager APIs.

## 🚨 **BLOCKING REQUIREMENT: COLLABORATIVE APPROVAL POLICY** 🚨

**🛑 MANDATORY STOP CHECKPOINTS - VIOLATION = IMMEDIATE FAILURE 🛑**

**BEFORE ANY ACTION, AI MUST:**

1. **🔍 ANALYZE REQUEST** - Understand what user wants
2. **📋 PRESENT OPTIONS** - List 2-3 specific approaches with pros/cons
3. **❓ ASK FOR APPROVAL** - Use EXACT phrase: "What would you like me to do next?"
4. **⏸️ WAIT FOR EXPLICIT APPROVAL** - No file changes, implementations, or tool usage without confirmation
5. **🔄 OFFER ALTERNATIVES** - Always provide multiple paths forward

**🚫 FORBIDDEN WITHOUT APPROVAL:**
- Creating or editing any files
- Running terminal commands  
- Making API calls or searches
- Implementing any solutions
- Generating code examples

**✅ REQUIRED APPROVAL TEMPLATE:**
```
## Analysis:
[Brief analysis of the request]

## Options:
**Option 1:** [Approach 1 - pros/cons]
**Option 2:** [Approach 2 - pros/cons] 
**Option 3:** [Alternative approach]

## What would you like me to do next?
```

**🚨 ENFORCEMENT TRIGGERS:**
- User requests implementation = STOP → Present options
- User asks questions = STOP → Present options  
- User mentions problems = STOP → Present options
- ANY request = STOP → Present options

**This takes ABSOLUTE PRIORITY over implementation speed.**

**Quick navigation:** [📚 Stack](#stack) | [🏗️ Project Structure](#project-structure) | [💬 Comment Policy](#🚨-blocking-enforcement-zero-tolerance-comment-policy-🚨) | [🚨 Comment Audit](#🚨-mandatory-comment-audit-checklist-🚨) | [🎯 Priority Enforcement](#🎯-ai-development-priority-enforcement) | [🚨 Testing Policy](#🚨-critical-testing-rule-policy-no-redundant-field-validation) | [🎯 AI Guidelines](#🎯-ai-development-guidelines) | [⚡ Implementation](#implementation-approaches) | [📖 Generic Guidelines](#generic-guidelines) | [🔍 API Discovery](#azure-api-discovery-process) | [🧠 Smart Patterns](#smart-pattern-recognition) | [❌ Error Handling](#error-handling-standards) | [🔧 Implementation Guide](#detailed-implementation-guidance) | [📚 Quick Reference](#quick-reference-links)

## Stack

- Go 1.22.x or later
- Terraform Plugin SDK
- Azure SDK for Go
- HashiCorp Go Azure SDK
- HashiCorp Go Azure Helpers
- Azure Resource Manager APIs

---
[⬆️ Back to top](#custom-instructions)

## Project Structure

```text
/internal
  /acceptance        - Acceptance test framework and helpers
  /clients           - Azure client configurations and setup
  /common            - Common utilities and helpers
  /features          - Feature flag management
  /provider          - Main provider configuration
  /sdk               - Internal SDK framework for typed resource implementations
  /services          - Service-specific implementations
    /servicename     - Individual Azure service (e.g., compute, storage)
      /client        - Service-specific client setup
      /parse         - Resource ID parsing utilities
      /validate      - Validation functions
      *.go           - Resource implementations
  /tf                - Terraform-specific utilities
/examples            - Terraform configuration examples
/website             - Provider documentation
/scripts             - Build and maintenance scripts
/vendor              - Go dependencies (managed by go mod)
```

---
[⬆️ Back to top](#custom-instructions)

## 🚨 BLOCKING ENFORCEMENT: ZERO TOLERANCE COMMENT POLICY 🚨

**🛑 BEFORE ANY CODE GENERATION - MANDATORY STOP**

**AI MUST ASK ITSELF THESE QUESTIONS BEFORE WRITING ANY CODE:**
1. "Am I about to add ANY comments (`//`, `/*`) to this code?"
2. "If YES → STOP and justify under 4-exception criteria"
3. "If NO exception applies → REFACTOR instead of commenting"

**ABSOLUTE RULE: NO UNNECESSARY COMMENTS**

Code must be self-documenting. Comments are the exception, not the rule.

**🚫 DEFAULT: Write code WITHOUT comments**

**Comments ONLY for these 4 cases:**
- Azure API-specific quirks not obvious from code
- Complex business logic that cannot be simplified
- Azure SDK workarounds for limitations/bugs
- Non-obvious state patterns (PATCH operations, residual state)

**🚫 NEVER comment these:**
- Variable assignments or struct initialization
- Standard Terraform/Go patterns
- Self-explanatory function calls
- Field mappings or obvious logic
- Error handling or nil checks

**3-SECOND RULE: Before ANY comment:**
1. Can I refactor instead? → **YES: Refactor, don't comment**
2. Is this an Azure API quirk? → **MAYBE: Comment acceptable**
3. Is this self-explanatory? → **YES: NO COMMENT**

**🔍 MANDATORY JUSTIFICATION:**
Every comment requires explicit justification:
- Which of the 4 exception cases applies?
- Why code cannot be self-explanatory?
- What specific Azure behavior needs documentation?

**FINAL CHECK:** "Can I eliminate this comment through better code?"

## 🚨 MANDATORY COMMENT AUDIT CHECKLIST 🚨

**BEFORE ANY CODE SUBMISSION - BLOCKING REQUIREMENT:**

**Step 1: AUTOMATIC VIOLATION SCAN**
- [ ] Search ALL code for comment patterns: `//`, `/*`, `#`
- [ ] For EACH comment found, answer: "Which of the 4 exception cases applies?"
  1. Azure API-specific quirks not obvious from code
  2. Complex business logic that cannot be simplified
  3. Azure SDK workarounds for limitations/bugs
  4. Non-obvious state patterns (PATCH operations, residual state)
- [ ] If NONE apply → **MANDATORY REMOVAL**

**Step 2: JUSTIFICATION DOCUMENTATION**
- [ ] For each remaining comment, document in review response:
  - Exception case number (1-4)
  - Why code cannot be self-explanatory through refactoring
  - Specific Azure behavior requiring documentation

**Step 3: REFACTORING VERIFICATION**
- [ ] Attempted better variable naming? **YES/NO**
- [ ] Attempted function extraction? **YES/NO**
- [ ] Attempted structure reorganization? **YES/NO**
- [ ] Cannot eliminate through code improvement? **YES/NO**

**🚫 FAILURE TO COMPLETE AUDIT = CODE REJECTION**

**FORBIDDEN COMMENT PATTERNS - Auto-Reject:**
```go
// Variable assignment comments
var enabled = true // Set enabled to true → DELETE

// Function call explanations
client.Create(params) // Create the resource → DELETE

// Obvious operations
d.Set("name", name) // Set the name field → DELETE

// Standard patterns
if err != nil { // Check for errors → DELETE
    return err
}
```

**ACCEPTABLE COMMENT PATTERNS:**
```go
// Azure Front Door PATCH operations preserve existing values when fields
// are omitted. Must explicitly send enabled=false to disable features.
result.Feature.Enabled = pointer.To(false)

// Azure SDK bug: nil values get filtered out before API call
// Workaround: always provide complete structure
```

**FORBIDDEN COMMENT PATTERNS - Auto-Reject:**
```go
// Variable assignment comments
var enabled = true // Set enabled to true → DELETE

// Function call explanations
client.Create(params) // Create the resource → DELETE

// Obvious operations
d.Set("name", name) // Set the name field → DELETE

// Standard patterns
if err != nil { // Check for errors → DELETE
    return err
}
```

**ACCEPTABLE COMMENT PATTERNS:**
```go
// Azure Front Door PATCH operations preserve existing values when fields
// are omitted. Must explicitly send enabled=false to disable features.
result.Feature.Enabled = pointer.To(false)

// Azure SDK bug: nil values get filtered out before API call
// Workaround: always provide complete structure
```

📋 **For comprehensive enforcement guidelines and detailed examples, see:** [Code Clarity Enforcement Guidelines](./instructions/code-clarity-enforcement.instructions.md#🚫-zero-tolerance-for-unnecessary-comments-policy)

---
[⬆️ Back to top](#custom-instructions)

## 🚨 CRITICAL TESTING RULE POLICY: No Redundant Field Validation

**MANDATORY: When using `data.ImportStep()`, DO NOT add redundant field validation checks**

```go
// ✅ CORRECT - Only ExistsInAzure check
func TestAccResourceName_basic(t *testing.T) {
    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r), // ONLY THIS
            ),
        },
        data.ImportStep(), // ImportStep validates ALL field values automatically
    })
}

// ❌ FORBIDDEN - Redundant field validation
func TestAccResourceName_basic(t *testing.T) {
    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("field").HasValue("value"), // FORBIDDEN
            ),
        },
        data.ImportStep(), // Already validates the field above
    })
}
```

**Exception**: Only add Key() checks for computed fields, TypeSet behavior, or Azure API transformations NOT in the configuration.

---
[⬆️ Back to top](#custom-instructions)

## 🎯 AI Development Guidelines

### Quick Decision Framework

**For ANY resource implementation task:**
1. **🔍 Analyze**: Determine if resource exists → Choose Typed (new) vs Untyped (maintenance)
2. **📋 Structure**: Follow checklist in [`implementation-guide.instructions.md`](./instructions/implementation-guide.instructions.md)
3. **⚡ Implement**: Use proven patterns from instruction files
4. **🧪 Test**: Create comprehensive acceptance tests
5. **📝 Document**: Follow templates in [`documentation-guidelines.instructions.md`](./instructions/documentation-guidelines.instructions.md)

### Code Quality Enforcement

**Before submitting ANY code:**
- ✅ All errors use `%+v` formatting with proper context
- ✅ Field names in error messages wrapped in backticks
- ✅ Azure SDK constants used for validation (not hardcoded values)
- ✅ Proper timeout configurations for Azure operations
- ✅ Import functionality tested and documented
- ✅ CustomizeDiff tested if validation logic exists
- ✅ No hardcoded values in tests or examples
- ✅ Resource ID parsing follows Azure patterns
- ✅ **CRITICAL: AI schema verification before field validation suggestions - AI must check Required/Optional/Optional+Computed before suggesting any empty/exists checks**
- ✅ **CRITICAL: Documentation fields follow alphabetical ordering - Required fields first (A-Z), then Optional fields (A-Z)**
- ✅ **CRITICAL: Code comments follow strict policy - only for Azure API quirks, complex business logic, or SDK workarounds**
- ✅ **All code comments have explicit justification documented in review response**
- ✅ **No comments on obvious operations, standard patterns, or self-explanatory code**

---
[⬆️ Back to top](#custom-instructions)

## Implementation Approaches

This provider supports two implementation approaches. **For comprehensive implementation patterns, detailed examples, and best practices, see the specialized instruction files in the repository.**

### Typed Resource Implementation (Preferred)
**Recommended for all new resources and data sources**

- **Framework**: Uses the internal/sdk framework with type-safe models
- **Model structures**: Struct types with tfschema tags for schema mapping
- **CRUD methods**: Receiver methods on resource struct types
- **State management**: metadata.Decode() and metadata.Encode() patterns
- **Client access**: metadata.Client for Azure SDK clients and structured logging
- **Error handling**: metadata.ResourceRequiresImport() and metadata.MarkAsGone()
- **Resource ID management**: metadata.SetID() for resource identification
- **Logging**: Structured logging through metadata.Logger

**Detailed Guidance**: See [`implementation-guide.instructions.md`](./instructions/implementation-guide.instructions.md) for comprehensive typed implementation patterns.

### Untyped Resource Implementation (Maintenance Only)
**Maintained for existing resources but not recommended for new development**

- **Framework**: Traditional Plugin SDK patterns
- **Function-based CRUD**: Functions like ResourceNameCreate, ResourceNameRead
- **State management**: Direct d.Set() and d.Get() patterns
- **Client access**: Traditional meta.(*clients.Client) initialization
- **Error handling**: tf.ImportAsExistsError() and direct state manipulation
- **Resource ID management**: Direct d.SetId() calls
- **Logging**: Traditional logging patterns

**Detailed Guidance**: See [`implementation-guide.instructions.md`](./instructions/implementation-guide.instructions.md) for untyped implementation maintenance patterns.

### Implementation-Aware Development
- **Code Review**: Both approaches must follow the same quality standards
- **Testing**: Identical acceptance test patterns regardless of implementation
- **Documentation**: User-facing documentation should be consistent between approaches
- **Azure Integration**: Both approaches integrate with the same Azure APIs and follow the same Azure-specific patterns

### CustomizeDiff Implementation Pattern

**IMPORTANT**: The dual import pattern is **only** required for specific scenarios:

**When DUAL IMPORTS are Required:**
```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"            // For *schema.ResourceDiff
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk" // For helpers
)

// When using *schema.ResourceDiff directly in CustomizeDiff functions
CustomizeDiff: pluginsdk.All(
    pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
        // Custom validation using *schema.ResourceDiff
        return nil
    }),
),
```

**When SINGLE IMPORT is Sufficient (Legacy Resources):**
```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"            // Only this import needed
)

// When using *pluginsdk.ResourceDiff in CustomizeDiffShim functions
CustomizeDiff: pluginsdk.CustomDiffWithAll(
    pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
        // Custom validation using *pluginsdk.ResourceDiff (which is an alias for *schema.ResourceDiff)
        return nil
    }),
),
```

**Rule of Thumb:**
- **Typed Resources**: Usually need dual imports when using `*schema.ResourceDiff` directly
- **Legacy/Untyped Resources**: Usually only need schema import when using `*pluginsdk.ResourceDiff`
- **Check the function signature**: If you see `*pluginsdk.ResourceDiff`, single import is sufficient

**Standard CustomizeDiff Resource Pattern:**
```go
func resourceServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceServiceNameCreate,
        Read:   resourceServiceNameRead,
        Update: resourceServiceNameUpdate,
        Delete: resourceServiceNameDelete,

        CustomizeDiff: pluginsdk.All(
            // Must use *schema.ResourceDiff from external package
            pluginsdk.ForceNewIfChange("property_name", func(ctx context.Context, old, new, meta interface{}) bool {
                return old.(string) != new.(string)
            }),
            func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
                // Custom validation logic
                if diff.Get("enabled").(bool) && diff.Get("configuration") == nil {
                    return fmt.Errorf("configuration is required when enabled is true")
                }
                return nil
            },
        ),

        Schema: map[string]*pluginsdk.Schema{
            // Schema definitions use pluginsdk types
        },
    }
}
```

**Why This Pattern is Required:**
- The internal pluginsdk package provides aliases for most Plugin SDK types
- However, CustomizeDiff function signatures are **not** aliased and must use *schema.ResourceDiff
- The pluginsdk.All(), pluginsdk.ForceNewIfChange() helpers are available in the internal package
- Resource and schema definitions use pluginsdk types for consistency

---
[⬆️ Back to top](#custom-instructions)

## Generic Guidelines

### Resource Implementation Guidelines
- Follow the standard resource lifecycle: Create, Read, Update, Delete (CRUD)
- Use proper Terraform Plugin SDK patterns and conventions
- Implement proper state management and drift detection
- Use ForceNew for properties that require resource recreation
- Implement proper timeout configurations
- Use appropriate validation functions for resource properties

### Azure API Integration
- Use the official Azure SDK for Go when available
- Implement proper error handling for Azure API responses
- Use appropriate polling for long-running operations
- Implement proper retry logic with exponential backoff
- Handle Azure API rate limits and throttling

### Azure API Discovery Process
**Efficient discovery of Azure APIs for provider implementation**

**🚀 PRIMARY DISCOVERY METHOD (Use First - Maximum 2 Minutes)**

**Step 1: Direct SDK Repository Search**
```bash
# Use github_repo tool to search HashiCorp Go Azure SDK directly
github_repo tool: "hashicorp/go-azure-sdk" 
Query: "{service-name} {resource-type} {api-version}"
# Example: "hardwaresecuritymodules cloudhsmclusters 2025-03-31"
```

**Step 2: Vendor Folder Verification**
```bash
# Check if API already exists in local vendor directory
grep_search: "vendor/**" pattern: "{service}/{version}/{resource}"
# Example: "hardwaresecuritymodules/2025-03-31/cloudhsmclusters"
```

**⏱️ Time Guideline: If not found in 2 minutes → API likely doesn't exist or needs custom implementation**

**✅ Success Criteria:**
- Found model files (model_*.go) with resource schema
- Found method files (method_*.go) with CRUD operations
- Found client.go with API client setup
- Confirmed import path: `github.com/hashicorp/go-azure-sdk/resource-manager/{service}/{version}/{resource}`

---

**📚 REFERENCE ONLY (Detailed Investigation When Primary Method Fails)**

**🔍 Repository Location**: `https://github.com/hashicorp/go-azure-sdk/tree/main/resource-manager`

**📁 Structure Pattern**:
```text
resource-manager/
├── {service-name}/              # e.g., hardwaresecuritymodules, compute, storage
│   ├── {api-version}/           # e.g., 2025-03-31, 2024-01-01
│   │   ├── {resource-type}/     # e.g., cloudhsmclusters, virtualmachines
│   │   │   ├── README.md        # Usage examples and documentation
│   │   │   ├── client.go        # Client initialization
│   │   │   ├── model_*.go       # Data models and types
│   │   │   ├── method_*.go      # CRUD operations
│   │   │   └── id_*.go          # Resource ID structures
```

**⚡ Manual Discovery Steps (Reference Only)**:
1. **Navigate to service**: Find folder matching Azure service (e.g., `hardwaresecuritymodules`)
2. **Choose API version**: Select latest available version (e.g., `2025-03-31`)
3. **Locate resource type**: Find specific resource folder (e.g., `cloudhsmclusters`)
4. **Review README.md**: Check available operations and usage examples
5. **Verify operations**: Confirm CRUD methods exist (`CreateOrUpdate`, `Get`, `Update`, `Delete`)

**✅ Standard Operations to Verify**:
- `CreateOrUpdate` / `CreateOrUpdateThenPoll` - Resource creation
- `Get` - Resource retrieval  
- `Update` / `UpdateThenPoll` - Resource updates (if supported)
- `Delete` / `DeleteThenPoll` - Resource deletion
- `ListByResourceGroup` - List resources by resource group
- `ListBySubscription` - List resources by subscription (if applicable)

**📋 Key Files to Check**:
- `client.go` - Client initialization patterns
- `model_*.go` - Resource data structures and properties
- `method_*.go` - Available operations and HTTP methods
- `id_*.go` - Resource ID parsing and construction

**🔗 Import Path Pattern**:
```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/{service}/{version}/{resource}"
```

**Example - Cloud HSM Cluster**:
```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2025-03-31/cloudhsmclusters"

// Client operations available:
client.CreateOrUpdateThenPoll(ctx, id, model)
client.Get(ctx, id)
client.UpdateThenPoll(ctx, id, patchModel)  
client.DeleteThenPoll(ctx, id)
client.ListByResourceGroup(ctx, resourceGroupId)
```

**🚨 Version Selection Strategy**:
- ✅ **Latest Version**: Use most recent API version for new implementations
- ⚠️ **Preview Versions**: Avoid preview APIs for production resources
- 📝 **Compatibility**: Ensure version supports all required operations

**❌ What to Do if API Not Found**:
1. Check if resource is available in a different service folder
2. Verify API version availability (might be in older/newer version)
3. Search for similar resource types for reference patterns
4. Check Azure REST API documentation for correct service categorization
5. Consider if resource requires custom SDK implementation

### Go Module Workflow for New Azure APIs

**🔧 Standard Workflow** (Based on Provider Standards):

When adding support for new Azure APIs, follow this proven workflow:

```bash
# 1. FIRST: Add the import to your Go code
import "github.com/hashicorp/go-azure-sdk/resource-manager/{service}/{version}/{resource}"

# 2. THEN: Let Go automatically resolve and fetch dependencies
go mod tidy

# 3. FINALLY: Update the vendor directory
go mod vendor

# 4. Test your implementation
make test
```

**⚡ Why This Order Works**:
- Go automatically detects the new import and fetches the required version
- `go mod tidy` resolves dependency conflicts and updates `go.mod`/`go.sum`
- `go mod vendor` copies dependencies to the `/vendor` directory
- The provider uses vendored dependencies for reproducible builds

**🚨 Version Selection Strategy**:
- **User Specified Version**: If user requests specific API version (e.g., "2023-01-01"), use that exact version
- **Latest Version (Default)**: If no version specified, use latest available and inform user:
  ```
  "Using latest API version 2025-03-31 for {service}. Specify a different version if needed."
  ```

**📋 Based on Provider Tools**:
This workflow is used by the provider's automated tools:
- `scripts/update-deps.sh` - Dependency management script
- `internal/tools/update-go-azure-sdk` - SDK version update automation
- `GNUmakefile` - Build verification process

**Example Implementation Flow**:
```go
// 1. Add import for Cloud HSM Clusters API
import "github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2025-03-31/cloudhsmclusters"

// 2. Implement resource using the SDK
func (r CloudHsmClusterResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.HardwareSecurityModules.CloudHsmClustersClient
            // Implementation continues...
        },
    }
}

// 3. Run go mod tidy && go mod vendor
// 4. Test with make test
```

### Security Guidelines
- Never hardcode sensitive values in tests or examples
- Use Azure Key Vault references where appropriate
- Implement proper authentication handling
- Validate all inputs to prevent injection attacks
- Follow principle of least privilege for service principals

### Testing Guidelines
**For comprehensive testing patterns and implementation-specific guidance, see [`testing-guidelines.instructions.md`](./instructions/testing-guidelines.instructions.md)**
- Write comprehensive acceptance tests for all resources
- Use the standard acceptance test framework
- Mock external dependencies appropriately
- Test both success and failure scenarios
- Ensure tests are idempotent and can run in parallel
- Test patterns should be consistent regardless of implementation approach

### Performance Considerations
- Implement efficient resource queries
- Use bulk operations where supported by Azure APIs
- Implement proper caching where appropriate
- Monitor and optimize API call patterns
- Use context with appropriate timeouts

---
[⬆️ Back to top](#custom-instructions)

## 🎯 AI Development Priority Enforcement

### **BLOCKING PRIORITY #1: ZERO TOLERANCE COMMENT POLICY**

**MANDATORY ACTIONS FOR EVERY CODE OPERATION:**

🚫 **BEFORE ANY CODE GENERATION/REVIEW:**
1. **AUTOMATIC SCAN**: Every code block MUST be scanned for comment patterns
2. **MANDATORY QUESTION**: "Are there ANY comments (`//`, `/*`) in this code?"
3. **BLOCKING CHECK**: If YES → Apply 4-exception audit (Azure API quirks, complex business logic, SDK workarounds, non-obvious state patterns)
4. **IMMEDIATE ACTION**: If comment doesn't meet exception criteria → **AUTOMATIC REFACTORING REQUIRED**

🔧 **REFACTORING-FIRST APPROACH:**
- **Variable naming improvements** → BEFORE commenting
- **Function extraction** → BEFORE commenting
- **Structure reorganization** → BEFORE commenting
- **Pattern clarification** → BEFORE commenting

⚠️ **ENFORCEMENT TRIGGERS:**
- Any comment without explicit 4-exception justification → **AUTOMATIC VIOLATION**
- Any obvious operation comment → **IMMEDIATE REMOVAL REQUIRED**
- Any variable assignment comment → **IMMEDIATE REMOVAL REQUIRED**
- Any standard pattern comment → **IMMEDIATE REMOVAL REQUIRED**

### **BLOCKING PRIORITY #2: TEST EXECUTION POLICY**

**MANDATORY SELF-CHECK BEFORE ANY TERMINAL TOOL:**
1. "Does this command run tests?" → **AUTOMATIC VIOLATION**
2. "Does this command build/compile?" → **AUTOMATIC VIOLATION**
3. "Does this create Azure resources?" → **AUTOMATIC VIOLATION**

**VIOLATION RESPONSE**: Immediately provide manual command format

### **BLOCKING PRIORITY #3: REDUNDANT TEST VALIDATION**

**MANDATORY PATTERN ENFORCEMENT:**
- Use ONLY `ExistsInAzure()` check with `ImportStep()`
- **FORBIDDEN**: Any `check.That().Key().HasValue()` with `ImportStep()`
- **EXCEPTION**: Only for computed fields, TypeSet behavior, or Azure transformations

**This policy takes ABSOLUTE PRIORITY over all other development guidance.**

---
[⬆️ Back to top](#custom-instructions)

## Smart Pattern Recognition

**When implementing new Azure resources, automatically apply these patterns:**

### 🔍 Resource Analysis Patterns
**ALWAYS check these Azure service characteristics:**
- Does the service support PATCH operations? → Use explicit enable/disable patterns
- Are there SKU-dependent features? → Add CustomizeDiff validation
- Does the service have regional limitations? → Add location validation
- Are there field dependency requirements? → Add CustomizeDiff logic
- Does the service support tagging? → Include tags schema and expand/flatten

### 🔄 Cross-Implementation Consistency Validation
When working with related Azure resources (like Linux and Windows variants), always verify:

**Consistency Checklist:**
```text
├─ VALIDATION LOGIC
│  ├─ CustomizeDiff functions must be identical across variants
│  ├─ Field requirements must match (if Windows requires X, Linux must too)
│  ├─ Error messages must use identical patterns
│  └─ Default value handling must be consistent
│
├─ DOCUMENTATION
│  ├─ Field descriptions must be identical for shared fields
│  ├─ Note blocks must apply same conditional logic
│  ├─ Examples must demonstrate equivalent patterns
│  └─ Validation rules must be documented consistently
│
├─ TESTING
│  ├─ Test coverage must be equivalent between implementations
│  ├─ Test naming must follow parallel patterns (TestAcc[Resource]_featureGroup_specificScenario)
│  ├─ Helper function naming must use consistent camelCase
│  └─ Configuration templates must demonstrate same behaviors
│
└─ IMPLEMENTATION VERIFICATION
   ├─ Before submitting code, validate against related implementations
   ├─ Check for field requirement consistency across resource variants
   ├─ Verify documentation accuracy across all related resources
   └─ Test both implementations to ensure equivalent behavior
```

### ⚡ Quick Implementation Patterns
```go
// NEW RESOURCE: Always start with this template
func (r ServiceNameResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceClient
            subscriptionId := metadata.Client.Account.SubscriptionId

            var model ServiceNameModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            id := parse.NewServiceNameID(subscriptionId, model.ResourceGroup, model.Name)

            metadata.Logger.Infof("Import check for %s", id)
            existing, err := client.Get(ctx, id)
            if err != nil && !response.WasNotFound(existing.HttpResponse) {
                return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
            }

            if !response.WasNotFound(existing.HttpResponse) {
                return metadata.ResourceRequiresImport(r.ResourceType(), id)
            }

            // Implementation continues...
            return nil
        },
    }
}
```

### 🧪 Testing Pattern Recognition

**ALWAYS include these tests for ANY resource:**

- **TestAcc[ResourceName]_basic** - Core functionality
- **TestAcc[ResourceName]_requiresImport** - Import conflict detection
- **TestAcc[ResourceName]_update** - If resource supports updates
- **TestAcc[ResourceName]_complete** - Comprehensive feature demonstration with all sensible field values
- **TestAcc[ResourceName]_customizeDiffValidation** - If CustomizeDiff exists

### Documentation
**For comprehensive documentation standards, see [`documentation-guidelines.instructions.md`](./instructions/documentation-guidelines.instructions.md)**
- Follow Terraform documentation standards
- Include comprehensive examples for all resources
- Document all resource attributes and their behaviors
- Include import documentation for existing resources
- Keep documentation synchronized with code changes
- Keep documentation up-to-date with code changes

### Git Workflow
- Use meaningful commit messages following conventional commits
- Create feature branches from main
- Squash commits before merging PRs
- Include issue numbers in commit messages
- Follow the contributor guidelines in /contributing directory

### Commit Message Guidelines
When writing commit messages, follow these standards for consistency and clarity:

**Format**: Use clear commit format: `TYPE: description`
- **Types**: ENHANCEMENT, BUG, DOCS, STYLE, REFACTOR, TEST, CHORE
- **Description**: Clear, concise summary under 50 characters

**For complex changes, include details:**
1. **Breaking Changes**: Specify if the change introduces breaking changes and describe the impact
2. **New Features**: Highlight any new resources, data sources, or significant features added
3. **Azure Services**: Mention updates to Azure services, APIs, or SDK integrations
4. **Line Length**: Keep each line under 72 characters for readability

**Examples**:
```text
ENHANCEMENT: add front door profile log scrubbing support

- Add log_scrubbing_rule schema for CDN Front Door profiles
- Implement expand/flatten functions for log scrubbing configuration
- Add comprehensive acceptance tests for log scrubbing validation
- Update documentation with log scrubbing examples

Closes Issue: #12345
```

```text
BUG: fix storage account network rules state drift

- Correct flatten function to handle empty network rules properly
- Fix state inconsistency when Azure returns default values
- Add acceptance test for network rules edge cases
- Ensure import functionality preserves correct state

Fixes Issue: #54321
```

```text
BREAKING CHANGE: remove deprecated storage account properties

- Remove deprecated `enable_blob_encryption` and `enable_file_encryption` fields
- These fields have been superseded by `encryption` block since v2.0
- Existing configurations using deprecated fields will fail validation
- Users must migrate to `encryption` block before upgrading provider
- See v4.0 upgrade guide for migration steps

Closes Issue: #87654
```

```text
DOCS: add explicit warning to fieldName for clarity

- Update documentation with clearer field usage warnings
- Add note blocks for conditional behavior requirements
- Improve user guidance for field combinations

Closes Issue: #98765
```

### Logging and Debugging
- Use structured logging with appropriate log levels
- Include resource IDs and operation context in logs
- Implement proper error wrapping and context
- Use Terraform's diagnostic system for user-facing errors
- Avoid logging sensitive information (passwords, keys, etc.)

### Go Code Standards
**For comprehensive coding standards, see [`implementation-guide.instructions.md`](./instructions/implementation-guide.instructions.md)**
- Follow effective Go practices and idioms
- Use gofmt for code formatting
- Implement proper error handling (don't ignore errors)
- Use context.Context for cancellation and timeouts
- Follow Go naming conventions (exported vs unexported)
- Use interfaces where appropriate for testing and modularity

### Terraform Provider Patterns
- Use pluginsdk.Resource for all resource definitions
- Implement proper schema validation
- Use CustomizeDiff for complex validation logic
- Implement proper ImportState functions
- Use appropriate TypeSet, TypeList, and TypeMap patterns
- Handle nested resource configurations properly

### Terraform Resource Lifecycle Patterns
- **Create**: Implement proper resource creation with all required properties
- **Read**: Always refresh state from Azure and handle deleted resources gracefully
- **Update**: Support in-place updates where possible, use ForceNew when necessary
- **Delete**: Handle cascading deletes and dependencies properly
- **Import**: Provide clear import documentation and test import functionality

### Azure-Specific Patterns
**For comprehensive Azure-specific guidance, see [`provider-guidelines.instructions.md`](./instructions/provider-guidelines.instructions.md)**
- Use standardized resource naming patterns with resourceToken
- Implement proper location/region handling
- Follow Azure resource tagging conventions
- Handle Azure API versioning correctly
- Implement proper subscription and resource group scoping
- Use Azure resource IDs consistently across resources

### Common Azure Resource Patterns
```go
// Standard resource schema pattern
func resourceExampleResource() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceExampleResourceCreate,
        Read:   resourceExampleResourceRead,
        Update: resourceExampleResourceUpdate,
        Delete: resourceExampleResourceDelete,

        Timeouts: &pluginsdk.ResourceTimeout{
            Create: pluginsdk.DefaultTimeout(30 * time.Minute),
            Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
            Update: pluginsdk.DefaultTimeout(30 * time.Minute),
            Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
        },

        Schema: map[string]*pluginsdk.Schema{
            "name": {
                Type:         pluginsdk.TypeString,
                Required:     true,
                ForceNew:     true,
                ValidateFunc: validation.StringIsNotEmpty,
            },
            "resource_group_name": commonschema.ResourceGroupName(),
            "location": commonschema.Location(),
            // ... other properties
        },
    }
}
```

### Development Workflow
1. **Setup**: Run `make tools` to install required tooling
2. **Build**: Use `make build` to compile the provider
3. **Test**: Run `make test` for unit tests, `make testacc` for acceptance tests
4. **Format**: Use `make fmt` to format code
5. **Lint**: Run `make lint` to check code quality
6. **Documentation**: Update examples and documentation as needed

### Acceptance Testing Guidelines
- Test both successful creation and error scenarios
- Test resource updates and partial updates
- Test resource import functionality
- Use appropriate test fixtures and cleanup
- Run tests against multiple Azure regions when relevant
- Use environment variables for test configuration (avoid hardcoding)
- Ensure test consistency regardless of implementation approach

### Code Review Checklist
- [ ] Follows Go coding standards and conventions
- [ ] Includes comprehensive acceptance tests
- [ ] **Tests use ONLY ExistsInAzure() check with ImportStep() - NO redundant field validation**
- [ ] Has proper error handling and logging
- [ ] Documentation is complete and accurate
- [ ] No hardcoded values or sensitive data
- [ ] Proper use of Terraform Plugin SDK patterns
- [ ] Azure API integration follows best practices
- [ ] Resource import functionality works correctly
- [ ] State management handles all edge cases
- [ ] Performance optimized for large deployments
- [ ] **CRITICAL: Code comments follow strict policy - only for Azure API quirks, complex business logic, or SDK workarounds**
- [ ] **All code comments have explicit justification documented in review response**
- [ ] **No comments on obvious operations, standard patterns, or self-explanatory code**

---
[⬆️ Back to top](#custom-instructions)

## Error Handling Standards

### Typed Resource Error Patterns
```go
// Use metadata.Decode for model decoding errors
var model ServiceNameResourceModel
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

### Untyped Error Patterns
```go
// Use consistent error formatting with context
if err != nil {
    return fmt.Errorf("creating Resource %q: %+v", name, err)
}

// Include resource information in error messages
if response.WasNotFound(resp.HttpResponse) {
    log.Printf("[DEBUG] Resource %q was not found - removing from state", id.ResourceName)
    d.SetId("")
    return nil
}

// Handle Azure-specific errors
if response.WasThrottled(resp.HttpResponse) {
    return resource.RetryableError(fmt.Errorf("request was throttled"))
}
```

### Common Error Standards (Both Approaches)
- Field names in error messages should be wrapped in backticks for clarity
- Field values in error messages should be wrapped in backticks for clarity
- Error messages must follow Go standards (lowercase, no punctuation, descriptive)
- Do not use contractions in error messages. Always use the full form of words. For example, write 'cannot' instead of 'can't' and 'is not' instead of 'isn't'
- Error messages must use '%+v' for verbose error output formatting
- Error messages must be clear, concise, and provide actionable guidance

---
[⬆️ Back to top](#custom-instructions)

## Detailed Implementation Guidance

### Choosing Implementation Approach
- **New Resources/Data Sources**: Always use Typed Resource Implementation
- **Existing Resources**: Continue using Untyped Resource Implementation for maintenance
- **Major Updates**: Consider migrating untyped resources to typed resource approach if significant changes are required
- **Bug Fixes**: Maintain existing implementation approach for simple bug fixes

### Typed Resource Implementation Best Practices
- Use type-safe model structures with appropriate `tfschema` tags
- Leverage receiver methods on resource struct types for CRUD operations
- Implement proper resource interfaces (`sdk.Resource`, `sdk.ResourceWithUpdate`, etc.)
- Use `metadata` for all client access, logging, and state management
- Follow structured error handling patterns with `metadata` methods
- Implement comprehensive validation in `IDValidationFunc()` method

### Untyped Resource Maintenance Best Practices
- Maintain existing function-based CRUD patterns
- Use direct schema manipulation with `d.Set()` and `d.Get()`
- Continue using traditional client initialization patterns
- Follow established error handling patterns with proper context
- Preserve existing resource behavior and state management
- Ensure backward compatibility when making changes

### Migration Considerations
- **User Experience**: Migration from untyped to typed resource should be transparent to users
- **State Compatibility**: Ensure Terraform state remains compatible across implementations
- **Feature Parity**: typed implementation should maintain all existing functionality
- **Testing**: Comprehensive testing required to validate migration behavior
- **Documentation**: Update internal development docs but keep user docs consistent

---
[⬆️ Back to top](#custom-instructions)

---

## Quick Reference Links

- 🏗️ **Implementation Guide**: [implementation-guide.instructions.md](./instructions/implementation-guide.instructions.md)
- 📋 **Code Clarity**: [code-clarity-enforcement.instructions.md](./instructions/code-clarity-enforcement.instructions.md)
- ⚡ **Azure Patterns**: [azure-patterns.instructions.md](./instructions/azure-patterns.instructions.md)
- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./instructions/testing-guidelines.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./instructions/documentation-guidelines.instructions.md)
- ❌ **Error Patterns**: [error-patterns.instructions.md](./instructions/error-patterns.instructions.md)
- 🔄 **Migration Guide**: [migration-guide.instructions.md](./instructions/migration-guide.instructions.md)
- 🏢 **Provider Guidelines**: [provider-guidelines.instructions.md](./instructions/provider-guidelines.instructions.md)
- 📐 **Schema Patterns**: [schema-patterns.instructions.md](./instructions/schema-patterns.instructions.md)

### 🚀 Enhanced Guidance Files

- ⚡ **Performance Optimization**: [performance-optimization.instructions.md](./instructions/performance-optimization.instructions.md)
- 🔐 **Security & Compliance**: [security-compliance.instructions.md](./instructions/security-compliance.instructions.md)
- 🔧 **Troubleshooting Decision Trees**: [troubleshooting-decision-trees.instructions.md](./instructions/troubleshooting-decision-trees.instructions.md)
- 🔄 **API Evolution Patterns**: [api-evolution-patterns.instructions.md](./instructions/api-evolution-patterns.instructions.md)

---
[⬆️ Back to top](#custom-instructions)
