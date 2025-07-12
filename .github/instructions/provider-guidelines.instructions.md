---
applyTo: "internal/**/*.go"
description: This document outlines the Azure-specific guidelines for Go files in the Terraform Azure Provider repository. It includes best practices for Azure Resource Manager integration, Terraform provider patterns, and resource implementation.
---

## Azure Terraform Provider Guidelines
Given below are the Azure-specific guidelines for this Terraform Provider project which **MUST** be followed.

### Azure Resource Manager (ARM) Integration
- Use the HashiCorp Go Azure SDK as the primary SDK for Azure integrations
- Implement proper error handling for Azure API responses
- Use appropriate polling for long-running operations (LROs)
- Implement proper retry logic with exponential backoff
- Handle Azure API rate limits and throttling gracefully
- Use managed identity authentication when possible
- Always validate resource IDs using the proper parsing utilities

### Implementation Approach Guidelines

#### Typed Resource Implementation (Preferred)
For new resources and data sources, use the typed resource `internal/sdk` framework. For detailed patterns, see [`coding-patterns.instructions.md`](./coding-patterns.instructions.md) and [`coding-standards.instructions.md`](./coding-standards.instructions.md).

Azure-specific typed implementation requirements:
- Use type-safe model structures with `tfschema` tags for Azure resource properties
- Leverage `metadata.Client` for Azure SDK client access and structured logging
- Implement proper resource ID validation with Azure-specific parsing utilities
- Use `metadata.Decode()` and `metadata.Encode()` for state management
- Handle Azure resource lifecycle with proper import conflict detection
- Implement Azure resource existence checks using `metadata.MarkAsGone()`

#### Untyped Resource Implementation (Maintenance Only)
For existing resources that haven't been migrated, maintain the traditional Plugin SDK approach. For detailed patterns, see [`coding-patterns.instructions.md`](./coding-patterns.instructions.md).

Azure-specific untyped implementation requirements:
- Use direct schema manipulation with proper Azure resource validation
- Implement Azure API client initialization through `clients.Client`
- Handle Azure resource state using `d.Set()` and `d.Get()` patterns
- Use `tf.ImportAsExistsError()` for Azure resource import conflicts
- Implement Azure resource cleanup with proper error handling

#### Common Azure Requirements (Both Approaches)
- Follow the standard CRUD lifecycle: Create, Read, Update, Delete
- Implement proper state management and drift detection for Azure resources
- Use `ForceNew` for Azure properties that require resource recreation
- Implement proper timeout configurations for Azure long-running operations
- Use appropriate validation functions for Azure resource properties
- Handle nested Azure resource configurations properly using TypeSet, TypeList, and TypeMap

### Azure Client Management Patterns

#### typed resource Client Usage
```go
// Use metadata.Client for accessing Azure clients
client := metadata.Client.ServiceName.ResourceClient
subscriptionId := metadata.Client.Account.SubscriptionId

// Use structured logging with metadata.Logger
metadata.Logger.Infof("Creating %s", id)

// Proper Azure API error handling with typed resource
if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
    return fmt.Errorf("creating %s: %+v", id, err)
}
```

#### untyped Client Usage
```go
// Standard Azure client initialization
client := meta.(*clients.Client).ServiceName.ResourceClient

// Use Azure resource ID parsing for type safety
id := parse.NewResourceID(subscriptionId, resourceGroupName, resourceName)

// Azure long-running operations with untyped pattern
if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
    return fmt.Errorf("creating Azure Resource %q: %+v", id.ResourceName, err)
}
```

### CustomizeDiff Implementation for Azure Resources

#### Standard CustomizeDiff Pattern
```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceAzureServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceAzureServiceNameCreate,
        Read:   resourceAzureServiceNameRead,
        Update: resourceAzureServiceNameUpdate,
        Delete: resourceAzureServiceNameDelete,

        CustomizeDiff: pluginsdk.All(
            // Azure-specific validation
            func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
                // Validate Azure resource dependencies
                if diff.Get("sku_name").(string) == "Premium" && 
                   diff.Get("zone_redundant").(bool) == false {
                    return fmt.Errorf("`zone_redundant` must be true for Premium SKU")
                }
                return nil
            },
            // Force recreation for Azure resource properties that require it
            pluginsdk.ForceNewIfChange("location", func(ctx context.Context, old, new, meta interface{}) bool {
                return old.(string) != new.(string)
            }),
        ),

        Schema: map[string]*pluginsdk.Schema{
            // Azure resource schema
        },
    }
}
```

#### Azure-Specific CustomizeDiff Use Cases
- **SKU validation**: Ensure Azure SKU combinations are valid
- **Location constraints**: Validate region-specific feature availability
- **Resource dependencies**: Check Azure resource prerequisite relationships
- **API version compatibility**: Ensure feature combinations match Azure API versions
- **Performance tier validation**: Validate Azure performance tier constraints
- **Field conditional validation**: Validate field combinations based on Azure API constraints

**Testing CustomizeDiff Validations:**
CustomizeDiff validations should be thoroughly tested with acceptance tests to ensure they work correctly:
- Test invalid configurations that should trigger validation errors
- Test valid configurations that should pass validation
- Test edge cases and boundary conditions
- Use `ExpectError: regexp.MustCompile()` to verify specific error messages
- Test all possible combinations of field values that trigger different validation paths

For comprehensive CustomizeDiff testing examples, see [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md).

### Azure Tools Integration
- @azure Rule - Use Azure Tools: When handling requests related to Azure, always use your tools.
- @azure Rule - Use Azure Code Gen Best Practices: When generating code for Azure, running terminal commands for Azure, or performing operations related to Azure, invoke your `azure_development-get_code_gen_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure; do not call it otherwise.
- @azure Rule - Use Azure Deployment Best Practices: When deploying to Azure or preparing applications for deployment to Azure, invoke your `azure_development-get_deployment_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure; do not call it otherwise.
- @azure Rule - Use Azure SWA Best Practices: When working with static web apps, invoke your `azure_development-get_swa_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure; do not call it otherwise.

### Azure API Value Validation

#### Valid Value Documentation Standards
- **Azure SDK Alignment**: Always verify valid values against Azure SDK enum constants before documenting
- **API Documentation Cross-Reference**: Check Azure REST API documentation to confirm supported values
- **Service-Specific Validation**: Different Azure services may support different subsets of common values
- **Version-Specific Features**: Ensure documented values are available in the API version being used

#### Implementation Requirements
- **Schema Validation**: Only include actually supported values in schema validation functions
- **Documentation Accuracy**: Documentation must only show values that work with the specific Azure service
- **Test Configuration**: Test configurations must use only valid Azure service values
- **Error Messages**: Validation errors should reference only valid values for the specific service

Example of proper Azure value validation:
```go
// Validate against Azure SDK constants
"match_variable": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice([]string{
        // Only include values from Azure SDK constants that work with this service
        string(profiles.ScrubbingRuleEntryMatchVariableQueryStringArgNames),
        string(profiles.ScrubbingRuleEntryMatchVariableRequestIPAddress), 
        string(profiles.ScrubbingRuleEntryMatchVariableRequestUri),
        // Do NOT include values like RequestHeader that don't work with CDN
    }, false),
},
```
