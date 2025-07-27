# Custom instructions

This is the official Terraform Provider for Azure (Resource Manager), written in Go. It enables Terraform to manage Azure resources through the Azure Resource Manager APIs.

## Stack

- Go 1.22.x or later
- Terraform Plugin SDK
- Azure SDK for Go
- HashiCorp Go Azure SDK
- HashiCorp Go Azure Helpers
- Azure Resource Manager APIs

## Project Structure

```
/internal
  /acceptance         - Acceptance test framework and helpers
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

## ⚠️ CRITICAL CODE COMMENT POLICY ⚠️

**NEVER add comments to code unless absolutely necessary.** Code should be self-documenting through clear variable names, function names, and structure.

**Comments should ONLY be added for:**
- Azure API-specific quirks or behaviors that are not obvious from code
- Complex business logic that cannot be made clear through code structure
- Workarounds for Azure SDK limitations or API bugs
- Non-obvious state management patterns (PATCH operations, residual state handling)
- Azure service constraints requiring explanation (timeout ranges, SKU limitations)

**DO NOT add comments for:**
- Variable assignments, struct initialization, basic operations
- Standard Terraform patterns (CRUD operations, schema definitions)
- Self-explanatory function calls or routine Azure API calls
- Field mappings between Terraform and Azure API models

**If the code is self-explanatory, do not add comments explaining what it does.**

## Implementation Approaches

This provider supports two implementation approaches. **For comprehensive implementation patterns, detailed examples, and best practices, see the specialized instruction files in the repository.**

### Typed Resource Implementation (Preferred)
**Recommended for all new resources and data sources**

- **Framework**: Uses the internal/sdk framework with type-safe models
- **Model structures**: Struct types with 	fschema tags for schema mapping
- **CRUD methods**: Receiver methods on resource struct types
- **State management**: metadata.Decode() and metadata.Encode() patterns
- **Client access**: metadata.Client for Azure SDK clients and structured logging
- **Error handling**: metadata.ResourceRequiresImport() and metadata.MarkAsGone()
- **Resource ID management**: metadata.SetID() for resource identification
- **Logging**: Structured logging through metadata.Logger

**Detailed Guidance**: See coding-patterns.instructions.md and coding-standards.instructions.md for comprehensive typed implementation patterns.

### Untyped Resource Implementation (Maintenance Only)
**Maintained for existing resources but not recommended for new development**

- **Framework**: Traditional Plugin SDK patterns
- **Function-based CRUD**: Functions like
esourceNameCreate,
esourceNameRead
- **State management**: Direct d.Set() and d.Get() patterns
- **Client access**: Traditional meta.(*clients.Client) initialization
- **Error handling**: 	f.ImportAsExistsError() and direct state manipulation
- **Resource ID management**: Direct d.SetId() calls
- **Logging**: Traditional logging patterns

**Detailed Guidance**: See coding-patterns.instructions.md for untyped implementation maintenance patterns.

### Implementation-Aware Development
- **Code Review**: Both approaches must follow the same quality standards
- **Testing**: Identical acceptance test patterns regardless of implementation
- **Documentation**: User-facing documentation should be consistent between approaches
- **Azure Integration**: Both approaches integrate with the same Azure APIs and follow the same Azure-specific patterns

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

### Security Guidelines
- Never hardcode sensitive values in tests or examples
- Use Azure Key Vault references where appropriate
- Implement proper authentication handling
- Validate all inputs to prevent injection attacks
- Follow principle of least privilege for service principals

### Testing Guidelines
**For comprehensive testing patterns and implementation-specific guidance, see `testing-guidelines.instructions.md`**
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

### Documentation
**For comprehensive documentation standards, see `documentation-guidelines.instructions.md`**
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

### Logging and Debugging
- Use structured logging with appropriate log levels
- Include resource IDs and operation context in logs
- Implement proper error wrapping and context
- Use Terraform's diagnostic system for user-facing errors
- Avoid logging sensitive information (passwords, keys, etc.)

### Go Code Standards
**For comprehensive coding standards, see `coding-standards.instructions.md` and `coding-style.instructions.md`**
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
**For comprehensive Azure-specific guidance, see `provider-guidelines.instructions.md`**
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
- [ ] Has proper error handling and logging
- [ ] Documentation is complete and accurate
- [ ] No hardcoded values or sensitive data
- [ ] Proper use of Terraform Plugin SDK patterns
- [ ] Azure API integration follows best practices
- [ ] Resource import functionality works correctly
- [ ] State management handles all edge cases
- [ ] Performance optimized for large deployments

### Azure Tools Integration
- @azure Rule - Use Azure Tools: When handling requests related to Azure, always use your tools.
- @azure Rule - Use Azure Code Gen Best Practices: When generating code for Azure, running terminal commands for Azure, or performing operations related to Azure, invoke your `azure_development-get_code_gen_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure; do not call it otherwise.
- @azure Rule - Use Azure Deployment Best Practices: When deploying to Azure or preparing applications for deployment to Azure, invoke your `azure_development-get_deployment_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure; do not call it otherwise.
- @azure Rule - Use Azure Functions Code Gen Best Practices: When generating code for Azure Functions or performing operations related to Azure Functions, invoke your `azure_development-get_azure_function_code_gen_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure Functions; do not call it otherwise.
- @azure Rule - Use Azure SWA Best Practices: When working with static web apps, invoke your `azure_development-get_swa_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure; do not call it otherwise.

### Error Handling Standards

#### Typed Resource Error Patterns
`go
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
`

#### Untyped Error Patterns
`go
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
`

#### Common Error Standards (Both Approaches)
- Field names in error messages should be wrapped in backticks for clarity
- Field values in error messages should be wrapped in backticks for clarity
- Error messages must follow Go standards (lowercase, no punctuation, descriptive)
- Do not use contractions in error messages. Always use the full form of words. For example, write 'cannot' instead of 'can't' and 'is not' instead of 'isn't'
- Error messages must use '%+v' for verbose error output formatting
- Error messages must be clear, concise, and provide actionable guidance

#### CustomizeDiff Implementation Pattern

**Dual Import Requirement:**
When implementing CustomizeDiff functions, both packages must be imported:

`go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)
`

**Standard CustomizeDiff Resource Pattern:**
`go
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
`

**Why This Pattern is Required:**
- The internal pluginsdk package provides aliases for most Plugin SDK types
- However, CustomizeDiff function signatures are **not** aliased and must use *schema.ResourceDiff
- The pluginsdk.All(), pluginsdk.ForceNewIfChange() helpers are available in the internal package
- Resource and schema definitions use pluginsdk types for consistency

### Detailed Implementation Guidance

#### Choosing Implementation Approach
- **New Resources/Data Sources**: Always use Typed Resource Implementation
- **Existing Resources**: Continue using Untyped Resource Implementation for maintenance
- **Major Updates**: Consider migrating untyped resources to typed resource approach if significant changes are required
- **Bug Fixes**: Maintain existing implementation approach for simple bug fixes

#### Typed Resource Implementation Best Practices
- Use type-safe model structures with appropriate `tfschema` tags
- Leverage receiver methods on resource struct types for CRUD operations
- Implement proper resource interfaces (`sdk.Resource`, `sdk.ResourceWithUpdate`, etc.)
- Use `metadata` for all client access, logging, and state management
- Follow structured error handling patterns with `metadata` methods
- Implement comprehensive validation in `IDValidationFunc()` method

#### Untyped Resource Maintenance Best Practices
- Maintain existing function-based CRUD patterns
- Use direct schema manipulation with `d.Set()` and `d.Get()`
- Continue using traditional client initialization patterns
- Follow established error handling patterns with proper context
- Preserve existing resource behavior and state management
- Ensure backward compatibility when making changes

#### Migration Considerations
- **User Experience**: Migration from untyped to typed resource should be transparent to users
- **State Compatibility**: Ensure Terraform state remains compatible across implementations
- **Feature Parity**: typed implementation should maintain all existing functionality
- **Testing**: Comprehensive testing required to validate migration behavior
- **Documentation**: Update internal development docs but keep user docs consistent
