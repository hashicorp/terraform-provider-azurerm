---
applyTo: "internal/website/docs/**/*.html.markdown"
description: This document outlines the standards and guidelines for writing documentation for Terraform resources and data sources in the AzureRM provider.
---

#  Documentation Guidelines

## Key Differences: Resources vs Data Sources

**Language Patterns:**
- **Resources**: Use action verbs - "Manages", "Creates", "Configures"
- **Data Sources**: Use retrieval verbs - "Gets information about", "Use this data source to access information about"

**Argument Types:**
- **Resources**: Arguments are for configuration (Required/Optional)
- **Data Sources**: Arguments are for identification/filtering (Required for lookup)

**Attributes:**
- **Resources**: Exports computed values after creation/update
- **Data Sources**: Exports all available information from existing resources

## Documentation Structure

**File Organization:**
```
website/docs/
 r/                          # Resource documentation
    service_resource.html.markdown
 d/                          # Data source documentation
    service_resource.html.markdown
 guides/                     # Provider guides and tutorials
     guide_name.html.markdown
```

**File Naming:**
- **Resources**: `r/service_resourcetype.html.markdown`
- **Data Sources**: `d/service_resourcetype.html.markdown`
- Use lowercase with underscores, match Terraform resource name exactly

## Resource Documentation Template

### Standard Resource Documentation Structure

```markdown
---
subcategory: "Service Name"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_resource"
description: |-
  Manages a Service Resource.
---

# azurerm_service_resource

Manages a Service Resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_resource" "example" {
  name                = "example-resource"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "Standard"

  tags = {
    environment = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Service Resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Service Resource should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Service Resource should exist. Changing this forces a new resource to be created.

* `sku_name` - (Required) The SKU name for this Service Resource. Possible values are `Standard` and `Premium`.

* `enabled` - (Optional) Whether this Service Resource is enabled. Defaults to `true`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Resource.

* `endpoint` - The endpoint URL of the Service Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Resource.
* `update` - (Defaults to 30 minutes) Used when updating the Service Resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Resource.

## Import

Service Resources can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_resource.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Service/resources/resource1
```

### Example Usage Subsections

**Standard Pattern**: Use only "## Example Usage" without subsections for most resources.

**When to Add Subsections**: Only create subsections under "Example Usage" when demonstrating meaningfully different configurations:

**Valid Subsection Scenarios:**
- **Platform variations**: "### Windows Function App" vs "### Linux Function App"
- **Authentication methods**: "### (with Base64 Certificate)" vs "### (with Key Vault Certificate)"
- **Deployment modes**: "### Standard Deployment" vs "### Premium Deployment with Custom Domain"
- **Integration patterns**: "### With Virtual Network" vs "### With Private Endpoint"

**Subsection Naming Convention:**
- Use descriptive names that clearly indicate the variation: "### Windows Function App"
- Include context in parentheses when helpful: "### (with Key Vault Certificate)"
- Avoid generic terms like "Basic", "Simple", or "Advanced"

**What NOT to create subsections for:**
- Minor field variations (add to main example instead)
- Single optional field demonstrations
- Tag variations or simple property changes
```

## Data Source Documentation Template

### Standard Data Source Documentation Structure

```markdown
---
subcategory: "Service Name"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_resource"
description: |-
  Gets information about an existing Service Resource.
---

# Data Source: azurerm_service_resource

Use this data source to access information about an existing Service Resource.

## Example Usage

```hcl
data "azurerm_service_resource" "example" {
  name                = "existing-resource"
  resource_group_name = "existing-resources"
}

output "service_resource_id" {
  value = data.azurerm_service_resource.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Service Resource.

* `resource_group_name` - (Required) The name of the Resource Group where the Service Resource exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Resource.

* `location` - The Azure Region where the Service Resource exists.

* `sku_name` - The SKU name of the Service Resource.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Service Resource.
```

## Writing Guidelines

### Language and Tone
- **Resources**: Use present tense action verbs ("manages", "creates", "configures")
- **Data Sources**: Use present tense retrieval verbs ("gets", "retrieves", "accesses")
- Be concise and clear, write for both beginners and experts

### Formatting Standards
- **Arguments**: Always use backticks around argument names: `argument_name`
- **Values**: Use backticks around specific values: `Standard`, `Premium`
- **Code blocks**: Use HCL syntax highlighting for Terraform code

### Front Matter Requirements
```yaml
---
subcategory: "Service Name"                                  # Azure service category
layout: "azurerm"                                            # Always "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_name"  # Page title
description: |-                                              # Brief description
  Manages a Service Resource.                                # For resources
  Gets information about an existing Service Resource.       # For data sources
---
```

**Note:** The `subcategory` value should match the service name from `./internal/services/[service-name]/registration.go`

### Resource Argument Documentation Patterns
```markdown
* `argument_name` - (Required) Description of what this argument configures. Additional context about behavior, constraints, or defaults.

* `complex_argument` - (Optional) A `complex_argument` block as defined below.

* `list_argument` - (Optional) A list of `list_argument` blocks as defined below.
```

### Data Source Argument Documentation Patterns
```markdown
* `argument_name` - The name/identifier used to locate the existing resource.

* `filter_argument` - The Filter criteria to narrow down the search results.
```

### Nested Block Documentation
```markdown
---

A `configuration` block supports the following:

* `required_setting` - (Required) Description of the required setting.

* `optional_setting` - (Optional) Description of the optional setting. Defaults to `default_value`.
```

## Example Configuration Guidelines

### The "None" Value Pattern in Documentation

When documenting resources that implement the "None" value pattern (where users omit optional fields instead of explicitly setting "None" values), examples should reflect this behavior:

**Example Considerations:**
- Show meaningful field access in outputs rather than fields that might be omitted due to the "None" pattern
- For log scrubbing examples, demonstrate accessing `match_variable` rather than `enabled` since `enabled` follows the "None" pattern
- Focus examples on fields that users actually configure and can reliably access

**Good Example Pattern:**
```hcl
output "log_scrubbing_match_variable" {
  value = data.azurerm_cdn_frontdoor_profile.example.log_scrubbing_rule.0.match_variable
}
```

**Pattern to Avoid:**
```hcl
output "log_scrubbing_enabled" {
  value = data.azurerm_cdn_frontdoor_profile.example.log_scrubbing_rule.0.enabled
}
```

The second example might not work as expected since `enabled` could be omitted when following the "None" value pattern.

### Example Configuration Strategy

When adding new fields to existing resources, follow this guidance for documentation examples:

**Update Existing Examples (Preferred for Simple Fields):**
- **Simple optional fields**: Add to existing basic/complete examples (e.g., `enabled = true`, `timeout_seconds = 300`)
- **Common configuration options**: Update existing examples rather than creating new ones (e.g., tags, basic settings)
- **Straightforward additions**: Fields that don't require complex explanation or setup

**Create New Examples (Only for Complex Features):**
- **Complex nested configurations**: Features requiring significant block structures or multiple related fields
- **Advanced use cases**: Features that require specific prerequisites or detailed explanation
- **Feature-specific scenarios**: When the field represents a distinct feature that warrants its own demonstration
- **Conditional configurations**: When field usage depends on specific combinations of other settings

**Example Decision Matrix:**
```markdown
 **Update existing**: `response_timeout_seconds = 120` (simple timeout field)
 **Update existing**: `enabled = false` (basic boolean toggle)
 **New example needed**: Complex nested `security_policy` block with multiple sub-configurations
 **New example needed**: Advanced `custom_domain` setup requiring certificates and DNS validation
```

This approach keeps documentation concise while ensuring complex features receive adequate explanation.

### Resource Example Requirements
- Include resource group creation
- Use realistic but generic names
- Show minimum required configuration
- Include common optional arguments (tags, location)
- Use consistent naming patterns
- Show creation and configuration patterns

### Data Source Example Requirements
- Show how to look up existing resources
- Demonstrate using data source outputs in other resources
- Include multiple lookup scenarios when relevant
- Show filtering and selection patterns
- Demonstrate practical use cases

### Resource Example Naming Conventions
```hcl
# Resource examples focus on creation
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_resource" "example" {
  name = "example-service"
  # Configuration for creation
}
```

### Data Source Example Naming Conventions
```hcl
# Data source examples focus on lookup
data "azurerm_service_resource" "example" {
  name                = "existing-service"
  resource_group_name = "existing-resources"
}

# Show how to use the data
output "service_endpoint" {
  value = data.azurerm_service_resource.example.endpoint
}

# Or use in other resources
resource "azurerm_other_resource" "example" {
  service_id = data.azurerm_service_resource.example.id
}
```

## Import Documentation

### Resource Import Format
````markdown
## Import

Service Resources can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_resource.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Service/resources/resource1
```
````

### Data Source Import (Not Applicable)
Data sources do not support import operations, so this section should be omitted from data source documentation.

## Timeout Documentation

### Resource Timeout Block
```markdown
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource.
* `update` - (Defaults to 30 minutes) Used when updating the Resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource.
```

### Data Source Timeout Block
```markdown
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Resource.
```

## Azure-Specific Documentation Patterns

### Resource Location Documentation
```markdown
* `location` - (Required) The Azure Region where the Resource should exist. Changing this forces a new resource to be created.
```

### Data Source Location Documentation
```markdown
* `location` - The Azure Region where the Resource exists.
```

### Resource Group Documentation
```markdown
# For Resources
* `resource_group_name` - (Required) The name of the Resource Group where the Resource should exist. Changing this forces a new resource to be created.

# For Data Sources
* `resource_group_name` - (Required) The name of the Resource Group where the Resource exists.
```

### Tags Documentation
```markdown
# For Resources
* `tags` - (Optional) A mapping of tags to assign to the resource.

# For Data Sources
* `tags` - A mapping of tags assigned to the resource.
```

### SKU Documentation
```markdown
# For Resources
* `sku_name` - (Required) The SKU name for this Resource. Possible values are `Standard_S1`, `Standard_S2`, and `Premium_P1`.

# For Data Sources
* `sku_name` - The SKU name of the Resource.
```

## Attributes Reference Differences

### Resource Attributes
- Focus on what becomes available after creation
- Include computed values and system-generated properties
- Show what can be referenced by other resources

```markdown
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Resource.

* `endpoint` - The endpoint URL of the Service Resource.
```

### Data Source Attributes
- Include all available information from the existing resource
- Show comprehensive details that can be used elsewhere
- Focus on what information is retrieved

```markdown
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Resource.

* `location` - The Azure Region where the Service Resource exists.

* `sku_name` - The SKU name of the Service Resource.

* `enabled` - Whether the Service Resource is enabled.

* `configuration` - A `configuration` block as defined below.

* `endpoint` - The endpoint URL of the Service Resource.

* `tags` - A mapping of tags assigned to the resource.
```

## Field Documentation Rules

### Field Ordering Standards
- **Required fields first**: Always list required fields before optional fields in argument documentation
- **Alphabetical within category**: Within required and optional groups, list fields alphabetically
- **Consistent structure**: Maintain the same field ordering pattern across all block documentation

### Note Format Standards
- **Use note blocks for conditional behavior**: When field usage depends on other field values, use note format instead of inline descriptions
- **Note syntax**: Use `~> **Note:**` format for important behavioral information
- **Clear conditional logic**: Explain exactly when fields are used vs ignored
- **Separate concerns**: Keep the main field description simple, use notes for complex conditional behavior

Example of proper field documentation:
```markdown
* `match_variable` - (Required) The variable to be scrubbed from the logs. Possible values are `QueryStringArgNames`, `RequestIPAddress`, and `RequestUri`.

* `enabled` - (Optional) Is this scrubbing rule enabled? Defaults to `true`.

* `operator` - (Optional) The operator to use for matching. Currently only `EqualsAny` is supported. Defaults to `EqualsAny`.

* `selector` - (Optional) The name of the query string argument to be scrubbed.

~> **Note:** The `selector` field is required when the `match_variable` is set to `QueryStringArgNames` and cannot be set when the `match_variable` is `RequestIPAddress` or `RequestUri`.
```

### Azure-Specific Documentation Standards
- **Valid values only**: Only document values that are actually supported by the Azure service
- **API validation**: Verify all possible values against Azure SDK constants and API documentation
- **Cross-reference validation**: When implementing similar features across resources, ensure consistent value documentation
- **SDK alignment**: Match documentation values with Azure SDK enum constants where applicable

### Cross-Implementation Documentation Consistency

When documenting related Azure resources (like Linux and Windows VMSS), ensure consistency across implementations:

**Field Documentation Consistency:**
- **Identical descriptions**: Use the same field descriptions for shared functionality across resource variants
- **Consistent validation rules**: Document the same validation requirements for equivalent fields
- **Synchronized note blocks**: Apply identical conditional logic notes to both implementations
- **Cross-reference accuracy**: When updating one variant's documentation, verify and update the related variant

**Common Mistakes to Avoid:**
- **Inconsistent rank field requirements**: Ensure both Linux and Windows VMSS document identical rank field usage patterns
- **Mismatched default value claims**: Verify that default value documentation matches actual Azure SDK behavior
- **Divergent validation patterns**: Maintain identical validation logic documentation across related resources

**Documentation Validation Checklist:**
- [ ] Field requirements match between Linux and Windows variants
- [ ] Default value claims verified against Azure SDK behavior
- [ ] Note blocks use consistent conditional logic across implementations
- [ ] Examples demonstrate the same patterns for equivalent functionality

## Provider Documentation Standards (Note Formatting)

### Note Block Standards
All notes should follow the exact same format (`(->|~>|!>) **Note:**`) where level of importance is indicated through the different types of notes as documented below.

Breaking changes should not be included in resource documentation notes:
- Breaking changes in a minor version should be added to the top of the changelog
- Breaking changes in a major version should be added to the upgrade guide

### Informational Note (`-> **Note:**`)
Use informational note blocks when providing additional useful information, recommendations and/or tips to the user.

**Example - Additional information on supported values:**
```markdown
* `type` - (Required) The type. Possible values include `This`, `That`, and `Other`.

-> **Note:** More information on each of the supported types can be found in [type documentation](https://docs.microsoft.com/azure/service-name/)
```

### Warning Note (`~> **Note:**`)
Use warning note blocks when providing information that the user needs to avoid certain errors, however if these errors are encountered they should not break anything or cause irreversible changes.

**Example - Conditional argument requirements:**
```markdown
* `optional_argument_enabled` - (Optional) Is the optional argument enabled? Defaults to `false`.

* `optional_argument` - (Optional) An optional argument.

~> **Note:** The argument `optional_argument` is required when `optional_argument_enabled` is set to `true`.
```

### Caution Note (`!> **Note:**`)
Use caution note blocks when providing critical information on potential irreversible changes, data loss or other things that can negatively affect a user's environment.

**Example - Irreversible changes:**
```markdown
* `irreversible_argument_enabled` - (Optional) Is irreversible argument enabled? Defaults to `false`.

!> **Note:** The argument `irreversible_argument_enabled` cannot be disabled after being enabled.
```

### Note Formatting Guidelines
- **Consistent format**: Always use the exact syntax patterns shown above
- **Appropriate level**: Choose the right note type based on the severity and impact of the information
- **Clear messaging**: Provide actionable information that helps users avoid problems
- **Avoid overuse**: Use notes for important information, not obvious functionality
- **Reference linking**: Include links to external documentation when helpful
