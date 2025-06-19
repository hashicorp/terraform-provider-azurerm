---
applyTo: "internal/**/*.go"
description: This document outlines the coding style for Go files in the Terraform AzureRM provider repository. It includes naming conventions, formatting rules, and guidelines for writing maintainable Terraform provider code.
---

## Coding Style for Go Files
Given below are the coding style guidelines for the Terraform AzureRM provider which **MUST** be followed.

### Copyright Header (Required)
All Go files must include this exact copyright header at the top:
```go
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
```

### Code Formatting (gofmt/gofumpt Enforced)
- **gofmt**: All Go code must be formatted with `gofmt` (automatically handled by most editors)
- **gofumpt**: Use `gofumpt` for stricter formatting that enforces additional style rules beyond gofmt
- **goimports**: Use `goimports` to automatically organize import statements
- **Indentation**: Use tabs (Go standard, handled by gofmt/gofumpt)
- **Line Length**: Aim for 120 characters max, break longer lines sensibly

### Import Organization
```go
package main

import (
    // Standard library imports first
    "context"
    "fmt"
    "log"
    
    // Third-party imports second
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
    
    // Local imports last
    "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
    "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
    "github.com/hashicorp/terraform-provider-azurerm/utils"
)
```

### Naming Conventions

#### Functions and Methods
- **Exported functions**: Use PascalCase (e.g., `CreateResource`, `ValidateInput`)
- **Unexported functions**: Use camelCase (e.g., `parseResourceID`, `buildParameters`)
- **Resource functions**: Follow pattern `resource[ServiceName][ResourceType][Operation]`
  - Examples: `resourceVirtualMachineCreate`, `resourceStorageAccountRead`

#### Variables and Constants
- **Exported variables**: Use PascalCase
- **Unexported variables**: Use camelCase
- **Constants**: Use PascalCase for exported, camelCase for unexported
- **Acronyms**: Keep as uppercase (e.g., `resourceGroupID`, `vmSSH`)

#### Types and Structs
- **Exported types**: Use PascalCase
- **Unexported types**: Use camelCase
- **Interface names**: Often end with 'er' (e.g., `ResourceProvider`, `Validator`)

### Terraform Provider Patterns

#### Resource Functions
```go
// Standard CRUD function signatures
func resourceServiceNameCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error
func resourceServiceNameRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error
func resourceServiceNameUpdate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error
func resourceServiceNameDelete(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error
```

#### Error Handling
```go
// Use consistent error formatting with context
if err != nil {
    return fmt.Errorf("creating Resource %q: %w", name, err)
}

// Include resource information in error messages
return fmt.Errorf("retrieving Resource %q (Resource Group %q): %w", 
    resourceName, resourceGroupName, err)
```

#### Client Usage
```go
// Standard client initialization pattern
client := meta.(*clients.Client).ServiceName.ResourceClient

// Use parsed resource IDs for type safety
id := parse.NewResourceID(subscriptionId, resourceGroupName, resourceName)
```

### Documentation and Comments

#### Function Documentation
```go
// resourceVirtualMachineCreate handles the creation of Azure Virtual Machines.
// It validates the configuration, provisions the VM through Azure Resource Manager,
// and waits for the deployment to complete before updating the Terraform state.
func resourceVirtualMachineCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
```

#### Package Documentation
```go
// Package compute provides Terraform resources for Azure Compute services
// including Virtual Machines, Virtual Machine Scale Sets, and related resources.
package compute
```

### Azure-Specific Patterns

#### Resource Schema
```go
// Use consistent schema patterns
"name": {
    Type:         pluginsdk.TypeString,
    Required:     true,
    ForceNew:     true,
    ValidateFunc: validation.StringIsNotEmpty,
},

"location": azure.SchemaLocation(),

"resource_group_name": azure.SchemaResourceGroupName(),
```

#### Long-Running Operations
```go
// Use *ThenPoll methods for long-running operations
if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
    return fmt.Errorf("creating Resource %q: %w", id.ResourceName, err)
}
```

### Code Quality Standards

#### golangci-lint Rules
- **gofmt**: Code must be formatted with gofmt
- **gofumpt**: Code must pass gofumpt stricter formatting checks
- **goimports**: Imports must be organized correctly
- **govet**: Pass go vet checks
- **ineffassign**: No ineffective assignments
- **misspell**: No common misspellings
- **unconvert**: No unnecessary type conversions

#### Testing Requirements
- All exported functions should have corresponding tests
- Use table-driven tests for multiple scenarios
- Mock external dependencies appropriately
- Follow acceptance test patterns for resource lifecycle testing

### File Organization

#### Directory Structure
- Place resource files in `internal/services/[service]/`
- Use separate files for different resource types
- Group related utilities in shared files
- Keep test files alongside source files with `_test.go` suffix

#### File Naming
- Resource files: `[service]_[resource_type]_resource.go`
- Data source files: `[service]_[resource_type]_data_source.go`
- Utility files: descriptive names (e.g., `validate.go`, `parse.go`)

### Build and Linting
- All code must pass `golangci-lint` checks
- Use `make lint` to run linting locally
- All tests must pass before merging
- Use `make test` for unit tests, `make testacc` for acceptance tests
