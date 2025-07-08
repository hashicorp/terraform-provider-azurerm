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


### Azure-Specific Resource Patterns
- Use standardized resource naming patterns with resourceToken
- Implement proper location/region handling across all resources
- Follow Azure resource tagging conventions consistently
- Handle Azure API versioning correctly
- Implement proper subscription and resource group scoping
- Use Azure resource IDs consistently across all resources
- Handle cascading deletes and dependencies properly

### Security and Authentication
- Never hardcode sensitive values in code, tests, or examples
- Use Azure Key Vault references where appropriate
- Implement proper authentication handling with service principals
- Validate all inputs to prevent injection attacks
- Follow principle of least privilege for service principals
- Use context with appropriate timeouts for all Azure API calls

### Testing Requirements
For comprehensive testing patterns and implementation details, see [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md).

Azure-specific testing requirements:
- Write comprehensive acceptance tests for all resources
- Test both success and failure scenarios thoroughly
- Ensure tests are idempotent and can run in parallel
- Test resource import functionality for all resources
- Use environment variables for test configuration (avoid hardcoding)
- Test Azure-specific features like resource tagging, location handling, and ARM template integration

### Error Handling and Logging
For implementation patterns and code examples, see [`coding-standards.instructions.md`](./coding-standards.instructions.md).

Azure-specific error handling requirements:
- Use structured logging with appropriate log levels
- Include Azure resource IDs and operation context in all logs
- Implement proper error wrapping and context propagation
- Handle Azure API errors gracefully with meaningful messages
- Avoid logging sensitive information (passwords, keys, tokens)
- Field names in error messages should be wrapped in backticks for clarity
- Field values in error messages should be wrapped in backticks for clarity
- Error messages must follow Go standards (lowercase, no punctuation, descriptive)
- Do not use contractions in error messages. Always use the full form of words. For example, write 'cannot' instead of 'can't' and 'is not' instead of 'isn't'
- Error messages must use '%+v' for verbose error output formatting
- Error messages must be clear, concise, and provide actionable guidance

### Performance Considerations
- Implement efficient resource queries to minimize API calls
- Use bulk operations where supported by Azure APIs
- Implement proper caching where appropriate
- Monitor and optimize API call patterns
- Use context with appropriate timeouts for all operations
- Consider pagination for large result sets

### Azure Tools Integration
- @azure Rule - Use Azure Tools: When handling requests related to Azure, always use your tools.
- @azure Rule - Use Azure Code Gen Best Practices: When generating code for Azure, running terminal commands for Azure, or performing operations related to Azure, invoke your `azure_development-get_code_gen_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure; do not call it otherwise.
- @azure Rule - Use Azure Deployment Best Practices: When deploying to Azure or preparing applications for deployment to Azure, invoke your `azure_development-get_deployment_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure; do not call it otherwise.
- @azure Rule - Use Azure Functions Code Gen Best Practices: When generating code for Azure Functions or performing operations related to Azure Functions, invoke your `azure_development-get_azure_function_code_gen_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure Functions; do not call it otherwise.
- @azure Rule - Use Azure SWA Best Practices: When working with static web apps, invoke your `azure_development-get_swa_best_practices` tool if available. Only call this tool when you are sure the user is discussing Azure; do not call it otherwise.
