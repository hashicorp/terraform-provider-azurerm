---
applyTo: "internal/**/*.go"
description: This document outlines the basic coding style for Go files in the Terraform AzureRM provider repository. It focuses on formatting rules, import organization, and fundamental Go style guidelines. For comprehensive standards and implementation patterns, see the other instruction files.
---

## Coding Style for Go Files
Given below are the coding style guidelines for the Terraform AzureRM provider which **MUST** be followed.

### Copyright Header (Required)
All Go files must include this exact copyright header at the top:
```go
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
```

### Implementation Guidelines
The following external contributing topics provide detailed best practices and guidance that supplement these coding style guidelines; new contributors are encouraged to review them for deeper context.
  - [`contributing/topics/best-practices.md`](../../contributing/topics/best-practices.md)
  - [`contributing/topics/guide-new-data-source.md`](../../contributing/topics/guide-new-data-source.md)
  - [`contributing/topics/guide-new-feature.md`](../../contributing/topics/guide-new-feature.md)
  - [`contributing/topics/guide-new-fields-to-data-source.md`](../../contributing/topics/guide-new-fields-to-data-source.md)
  - [`contributing/topics/guide-new-fields-to-resource.md`](../../contributing/topics/guide-new-fields-to-resource.md)
  - [`contributing/topics/guide-new-resource-vs-inline.md`](../../contributing/topics/guide-new-resource-vs-inline.md)
  - [`contributing/topics/guide-new-resource.md`](../../contributing/topics/guide-new-resource.md)
  - [`contributing/topics/guide-new-service-package.md`](../../contributing/topics/guide-new-service-package.md)
  - [`contributing/topics/guide-new-write-only-attribute.md`](../../contributing/topics/guide-new-write-only-attribute.md)
  - [`contributing/topics/guide-resource-ids.md`](../../contributing/topics/guide-resource-ids.md)
  - [`contributing/topics/guide-state-migrations.md`](../../contributing/topics/guide-state-migrations.md)
  - [`contributing/topics/reference-acceptance-testing.md`](../../contributing/topics/reference-acceptance-testing.md)
  - [`contributing/topics/reference-documentation-standards.md`](../../contributing/topics/reference-documentation-standards.md)
  - [`contributing/topics/reference-errors.md`](../../contributing/topics/reference-errors.md)
  - [`contributing/topics/schema-design-considerations.md`](../../contributing/topics/schema-design-considerations.md)
  - [`contributing/topics/reference-naming.md`](../../contributing/topics/reference-naming.md)

### Code Formatting (gofmt/gofumpt Enforced)
- **gofmt**: All Go code must be formatted with `gofmt` (automatically handled by most editors)
- **gofumpt**: Use `gofumpt` for stricter formatting that enforces additional style rules beyond gofmt
- **goimports**: Use `goimports` to automatically organize import statements
- **Indentation**: Use tabs (Go standard, handled by gofmt/gofumpt)
- **Line Length**: Aim for 120 characters max, break longer lines sensibly

### Import Organization

For comprehensive import organization patterns including Azure-specific CustomizeDiff requirements, see:
- **Detailed Patterns**: [coding-patterns.instructions.md](./coding-patterns.instructions.md) - Import Management Pattern section

**Quick Reference**: Follow three-group import structure:
1. **Standard library imports** (context, fmt, log, etc.)
2. **Third-party imports** (terraform-plugin-sdk, etc.)
3. **Local imports** (internal/clients, internal/services, etc.)

### Basic Go Naming Conventions

For comprehensive naming conventions including typed vs untyped patterns, see [`coding-standards.instructions.md`](./coding-standards.instructions.md).

#### Basic Rules
- **Exported identifiers**: Use PascalCase (e.g., `CreateResource`, `ValidateInput`)
- **Unexported identifiers**: Use camelCase (e.g., `parseResourceID`, `buildParameters`)
- **Acronyms**: Keep as uppercase (e.g., `resourceGroupID`, `vmSSH`, `apiURL`)
- **Interface names**: Often end with 'er' (e.g., `ResourceProvider`, `Validator`)

### typed vs untyped Style Considerations

While most formatting rules apply to both implementation approaches, there are some style differences to be aware of:

#### typed resource-Based Style (Preferred)
```go
// Model struct with tfschema tags (typed resource approach)
type ServiceResourceModel struct {
    Name          string            `tfschema:"name"`
    ResourceGroup string            `tfschema:"resource_group_name"`
    Location      string            `tfschema:"location"`
    Tags          map[string]string `tfschema:"tags"`
}

// Resource struct methods (typed resource approach)
func (r ServiceResourceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // Implementation
        },
    }
}
```

#### untyped Plugin SDK Style (Maintenance)
```go
// Traditional function-based approach (untyped)
func resourceServiceCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    // Implementation
}

// Traditional schema definition (untyped)
func resourceServiceSchema() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "name": {
            Type:     pluginsdk.TypeString,
            Required: true,
        },
    }
}
```

#### Style Guidelines for Both Approaches
- **File naming**: Use snake_case for all file names (`service_resource.go`, `service_data_source.go`)
- **Package organization**: Group related functionality within service packages
- **Import grouping**: Always maintain the three-group import structure (standard, third-party, local)
- **Function length**: Keep functions focused and under 50 lines when possible
- **Variable naming**: Use descriptive names that indicate purpose and scope

### Documentation and Comments

#### Go Documentation Standards
```go
// Package compute provides Terraform resources for Azure Compute services
// including Virtual Machines, Virtual Machine Scale Sets, and related resources.
package compute

// CreateResource creates a new Azure resource with the specified configuration.
// It validates the input parameters and returns an error if validation fails.
func CreateResource(config *ResourceConfig) error {
    // Implementation here
}
```

#### Documentation Style by Implementation Approach

**typed resource Documentation:**
```go
// ServiceResourceModel represents the Terraform model for Service Resource
type ServiceResourceModel struct {
    // Name is the unique identifier for the service resource
    Name string `tfschema:"name"`
}

// Create implements the create operation for Service Resource
func (r ServiceResourceResource) Create() sdk.ResourceFunc {
    // Return function handles resource creation lifecycle
}
```

**untyped Plugin SDK Documentation:**
```go
// resourceServiceCreate handles the creation of Service Resource in Terraform
func resourceServiceCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    // Extract configuration and create Azure resource
}

// resourceServiceSchema defines the Terraform schema for Service Resource
func resourceServiceSchema() map[string]*pluginsdk.Schema {
    // Return schema definition for resource configuration
}
```

#### Comment Style
- Use `//` for single-line comments
- Use `/* */` for multi-line comments only when necessary
- Comments should explain "why" not "what" when the code is clear
- Keep comments up-to-date with code changes

For detailed implementation patterns and comprehensive standards, see:
- [`coding-patterns.instructions.md`](./coding-patterns.instructions.md) - Detailed implementation patterns
- [`coding-standards.instructions.md`](./coding-standards.instructions.md) - Comprehensive coding standards
- [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md) - Testing guidelines
- [`provider-guidelines.instructions.md`](./provider-guidelines.instructions.md) - Azure-specific guidelines
