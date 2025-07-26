---
applyTo: "internal/website/docs/**/*.html.markdown"
description: This document outlines the standards and guidelines for writing documentation for Terraform resources and data sources in the AzureRM provider.
---

# Documentation Guidelines for Terraform AzureRM Provider

## Key Differences: Resources vs Data Sources

### Language Patterns
- **Resources**: Use action verbs - "Manages", "Creates", "Configures"
- **Data Sources**: Use retrieval verbs - "Gets information about", "Use this data source to access information about"

### Argument Types
- **Resources**: Arguments are for configuration (Required/Optional)
- **Data Sources**: Arguments are for identification/filtering (Required for lookup)

### Attributes
- **Resources**: Exports computed values after creation/update
- **Data Sources**: Exports all available information from existing resources

## Implementation Approach Considerations
## Implementation Approach Considerations

### Documentation Consistency Across Implementation Approaches

While the underlying Go implementation (typed resource vs untyped Plugin SDK) is transparent to end users, documentation should maintain consistency regardless of implementation approach.

**For detailed implementation approach information, see the main copilot instructions file.**

#### Key Documentation Standards
- **User Experience**: Documentation should be identical regardless of implementation approach
- **Feature Parity**: Both approaches should document the same Azure resource capabilities
- **Example Consistency**: HCL configuration examples should follow the same patterns
- **Behavioral Accuracy**: Document the resource behavior as users experience it

## Documentation Structure

### File Organization
```
website/docs/
├── r/                          # Resource documentation
│   └── service_resource.html.markdown
├── d/                          # Data source documentation
│   └── service_resource.html.markdown
└── guides/                     # Provider guides and tutorials
    └── guide_name.html.markdown
```

### File Naming Conventions
- **Resources**: `r/service_resourcetype.html.markdown`
  - Example: `r/cdn_frontdoor_profile.html.markdown`
- **Data Sources**: `d/service_resourcetype.html.markdown`
  - Example: `d/cdn_frontdoor_profile.html.markdown`
- Use lowercase with underscores for separation
- Match the Terraform resource/data source name exactly

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
```

# azurerm_service_resource

Manages a Service Resource.

## Example Usage

### Basic

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

### Complete

```hcl
# Include a more comprehensive example showing advanced features
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_resource" "example" {
  name                = "example-resource"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "Premium"
  enabled  = true

  configuration {
    setting1 = "value1"
    setting2 = "value2"
  }

  tags = {
    environment = "Production"
    project     = "Example"
  }
}
```

## Arguments Reference

````markdown
The following arguments are supported:

* `name` - (Required) The name of the Service Resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Service Resource should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Service Resource should exist. Changing this forces a new resource to be created.

* `sku_name` - (Required) The SKU name for this Service Resource. Possible values are `Standard` and `Premium`.

* `enabled` - (Optional) Whether this Service Resource is enabled. Defaults to `true`.

* `configuration` - (Optional) A `configuration` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `configuration` block supports the following:

* `setting1` - (Required) The first configuration setting.

* `setting2` - (Optional) The second configuration setting. Defaults to `default_value`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Resource.

* `endpoint` - The endpoint URL of the Service Resource.

* `status` - The current status of the Service Resource.

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
````

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
```

# Data Source: azurerm_service_resource

Use this data source to access information about an existing Service Resource.

## Example Usage

### Basic Lookup

```hcl
data "azurerm_service_resource" "example" {
  name                = "existing-resource"
  resource_group_name = "existing-resources"
}

output "service_resource_id" {
  value = data.azurerm_service_resource.example.id
}

output "service_resource_endpoint" {
  value = data.azurerm_service_resource.example.endpoint
}
```

### Using Data Source Output in Resources

```hcl
data "azurerm_service_resource" "example" {
  name                = "existing-resource"
  resource_group_name = "existing-resources"
}

resource "azurerm_other_resource" "example" {
  name                    = "example-other"
  service_resource_id     = data.azurerm_service_resource.example.id
  service_resource_endpoint = data.azurerm_service_resource.example.endpoint
}
```

## Arguments Reference

The following arguments are supported:

```markdown
* `name` - (Required) The name of this Service Resource.

* `resource_group_name` - (Required) The name of the Resource Group where the Service Resource exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Resource.

* `location` - The Azure Region where the Service Resource exists.

* `sku_name` - The SKU name of the Service Resource.

* `enabled` - Whether the Service Resource is enabled.

* `configuration` - A `configuration` block as defined below.

* `endpoint` - The endpoint URL of the Service Resource.

* `status` - The current status of the Service Resource.

* `tags` - A mapping of tags assigned to the resource.

---

A `configuration` block exports the following:

* `setting1` - The first configuration setting.

* `setting2` - The second configuration setting.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Service Resource.
```

## Language Guidelines

### Resource Documentation Language
- **Title**: Use action verbs - "Manages a...", "Creates a...", "Configures a..."
- **Description**: Focus on what the resource does
- **Examples**: Show configuration and creation patterns
- **Arguments**: Describe what values to provide for configuration
- **Attributes**: Describe what information is available after creation

### Data Source Documentation Language
- **Title**: Always start with "Data Source: " prefix
- **Description**: Use "Gets information about an existing..." or "Use this data source to access information about..."
- **Examples**: Show lookup patterns and how to use the retrieved data
- **Arguments**: Describe what values are needed to identify the resource
- **Attributes**: Describe what information is retrieved from the existing resource

## Documentation Standards

### Front Matter Requirements
```yaml
---
subcategory: "Service Name"                                  # Azure service category (e.g., "CDN", "Compute", "Network")
layout: "azurerm"                                            # Always "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_name"  # Page title for browser/search
description: |-                                              # Brief description (single line, with a period)
  Manages a Service Resource.                                # For resources
  Gets information about an existing Service Resource.       # For data sources
---
```

**Note:** The `subcategory` value should match the service name from the registration file:
* Navigate to `./internal/services/[service-name]/registration.go`
* Use the value returned by the `Name()` function
* **Example:** For CDN resources, check `./internal/services/cdn/registration.go` → `Name()` returns `"CDN"`
* **Example:** For Compute resources, check `./internal/services/compute/registration.go` → `Name()` returns `"Compute"`
* **Example:** For Network resources, check `./internal/services/network/registration.go` → `Name()` returns `"Network"`

The `subcategory` groups related resources in the Terraform Registry and provider documentation.

### Resource-Specific Front Matter
```yaml
---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_profile"
description: |-
  Manages a CDN Front Door Profile.
---
```

### Data Source-Specific Front Matter
```yaml
---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_profile"
description: |-
  Gets information about an existing CDN Front Door Profile.
---
```

### Subcategory Mapping
- **CDN**: Content Delivery Network resources
- **Compute**: Virtual Machines, Scale Sets, etc.
- **Containers**: Container services (AKS, ACI, etc.)
- **Database**: SQL, CosmosDB, MySQL, etc.
- **Network**: Virtual Networks, Load Balancers, etc.
- **Storage**: Storage Accounts, Blobs, etc.
- **Security**: Key Vault, Security Center, etc.
- **Web**: App Service, Function Apps, etc.

### Writing Guidelines

#### Language and Tone
- **Resources**: Use present tense action verbs ("manages", "creates", "configures")
- **Data Sources**: Use present tense retrieval verbs ("gets", "retrieves", "accesses")
- Be concise and clear
- Use active voice when possible
- Avoid unnecessary technical jargon
- Write for both beginners and experts

#### Formatting Standards
- **Arguments**: Always use backticks around argument names: `argument_name`
- **Values**: Use backticks around specific values: `Standard`, `Premium`
- **Code blocks**: Use HCL syntax highlighting for Terraform code
- **Lists**: Use asterisk (*) for unordered lists
- **Emphasis**: Use **bold** for important concepts, *italics* sparingly

#### Resource Argument Documentation Patterns
```markdown
* `argument_name` - (Required) Description of what this argument configures. Additional context about behavior, constraints, or defaults.

* `complex_argument` - (Optional) A `complex_argument` block as defined below.

* `list_argument` - (Optional) A list of `list_argument` blocks as defined below.
```

#### Data Source Argument Documentation Patterns
```markdown
* `argument_name` - The name/identifier used to locate the existing resource.

* `filter_argument` - The Filter criteria to narrow down the search results.
```

#### Nested Block Documentation
```markdown
---

A `configuration` block supports the following:

* `required_setting` - (Required) Description of the required setting.

* `optional_setting` - (Optional) Description of the optional setting. Defaults to `default_value`.
```

### Example Configuration Guidelines

#### The "None" Value Pattern in Documentation

When documenting resources that implement the "None" value pattern (where users omit optional fields instead of explicitly setting "None" values), examples should reflect this behavior:

**Example Considerations:**
- Show meaningful field access in outputs rather than fields that might be omitted due to the "None" pattern
- For log scrubbing examples, demonstrate accessing `match_variable` rather than `enabled` since `enabled` follows the "None" pattern
- Focus examples on fields that users actually configure and can reliably access

**Good Example Pattern:**
```hcl
output "log_scrubbing_match_variable" {
  value = data.azurerm_cdn_frontdoor_profile.example.log_scrubbing.0.scrubbing_rule.0.match_variable
}
```

**Pattern to Avoid:**
```hcl
output "log_scrubbing_enabled" {
  value = data.azurerm_cdn_frontdoor_profile.example.log_scrubbing.0.enabled
}
```

The second example might not work as expected since `enabled` could be omitted when following the "None" value pattern.

### Example Configuration Guidelines

#### The "None" Value Pattern in Documentation

When documenting resources that implement the "None" value pattern (where users omit optional fields instead of explicitly setting "None" values), examples should reflect this behavior:

**Example Considerations:**
- Show meaningful field access in outputs rather than fields that might be omitted due to the "None" pattern
- For log scrubbing examples, demonstrate accessing `match_variable` rather than `enabled` since `enabled` follows the "None" pattern
- Focus examples on fields that users actually configure and can reliably access

**Good Example Pattern:**
```hcl
output "log_scrubbing_match_variable" {
  value = data.azurerm_cdn_frontdoor_profile.example.log_scrubbing.0.scrubbing_rule.0.match_variable
}
```

**Pattern to Avoid:**
```hcl
output "log_scrubbing_enabled" {
  value = data.azurerm_cdn_frontdoor_profile.example.log_scrubbing.0.enabled
}
```

The second example might not work as expected since `enabled` could be omitted when following the "None" value pattern.

#### Resource Example Requirements
- Include resource group creation
- Use realistic but generic names
- Show minimum required configuration
- Include common optional arguments (tags, location)
- Use consistent naming patterns
- Show creation and configuration patterns

#### Data Source Example Requirements
- Show how to look up existing resources
- Demonstrate using data source outputs in other resources
- Include multiple lookup scenarios when relevant
- Show filtering and selection patterns
- Demonstrate practical use cases

#### Typed vs Untyped Implementation Examples
While implementation approach is transparent to users, ensure examples reflect actual behavior:

**typed implementation Considerations:**
- May have more comprehensive validation in examples
- Could include newer Azure features and capabilities
- May demonstrate improved error handling patterns
- Should reflect any enhanced functionality

**untyped Implementation Considerations:**
- Examples should match currently available functionality
- May have more basic validation patterns
- Should accurately represent current feature set
- Focus on proven, stable configuration patterns

**Consistency Requirements:**
- HCL syntax and structure should be identical regardless of implementation
- Resource and argument names must match exactly
- Behavior described should reflect actual user experience
- Examples should work with the current provider version

#### Resource Example Naming Conventions
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

#### Data Source Example Naming Conventions
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

### Import Documentation

#### Resource Import Format
````markdown
## Import

Service Resources can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_resource.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Service/resources/resource1
```
````

#### Data Source Import (Not Applicable)
Data sources do not support import operations, so this section should be omitted from data source documentation.

### Timeout Documentation

#### Resource Timeout Block
```markdown
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource.
* `update` - (Defaults to 30 minutes) Used when updating the Resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource.
```

#### Data Source Timeout Block
```markdown
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Resource.
```

### Azure-Specific Documentation Patterns

#### Resource Location Documentation
```markdown
* `location` - (Required) The Azure Region where the Resource should exist. Changing this forces a new resource to be created.
```

#### Data Source Location Documentation
```markdown
* `location` - The Azure Region where the Resource exists.
```

#### Resource Group Documentation
```markdown
# For Resources
* `resource_group_name` - (Required) The name of the Resource Group where the Resource should exist. Changing this forces a new resource to be created.

# For Data Sources
* `resource_group_name` - (Required) The name of the Resource Group where the Resource exists.
```

#### Tags Documentation
```markdown
# For Resources
* `tags` - (Optional) A mapping of tags to assign to the resource.

# For Data Sources
* `tags` - A mapping of tags assigned to the resource.
```

#### SKU Documentation
```markdown
# For Resources
* `sku_name` - (Required) The SKU name for this Resource. Possible values are `Standard_S1`, `Standard_S2`, and `Premium_P1`.

# For Data Sources
* `sku_name` - The SKU name of the Resource.
```

### Attributes Reference Differences

#### Resource Attributes
- Focus on what becomes available after creation
- Include computed values and system-generated properties
- Show what can be referenced by other resources

```markdown
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Resource.

* `endpoint` - The endpoint URL of the Service Resource.

```

#### Data Source Attributes
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

### Common Mistakes to Avoid

#### Resource Documentation Errors
- Using passive language instead of active management verbs
- Missing ForceNew behavior documentation
- Incorrect default values
- Missing import documentation
- Not showing practical configuration examples

#### Data Source Documentation Errors
- Forgetting "Data Source: " prefix in title
- Using creation language instead of retrieval language
- Including import documentation (not applicable)
- Not showing how to use the retrieved data
- Missing practical lookup examples

### Field Documentation Rules

#### Field Ordering Standards
- **Required fields first**: Always list required fields before optional fields in argument documentation
- **Alphabetical within category**: Within required and optional groups, list fields alphabetically
- **Consistent structure**: Maintain the same field ordering pattern across all block documentation

#### Note Format Standards
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

#### Azure-Specific Documentation Standards
- **Valid values only**: Only document values that are actually supported by the Azure service
- **API validation**: Verify all possible values against Azure SDK constants and API documentation
- **Cross-reference validation**: When implementing similar features across resources, ensure consistent value documentation
- **SDK alignment**: Match documentation values with Azure SDK enum constants where applicable

### Quality Checklist

#### Resource Documentation Review
- [ ] Title uses active management language
- [ ] Examples show creation and configuration
- [ ] All arguments documented with correct optionality
- [ ] ForceNew behavior documented where applicable
- [ ] Import documentation included
- [ ] Timeouts documented for all operations
- [ ] Practical configuration examples provided

#### Data Source Documentation Review
- [ ] Title starts with "Data Source: "
- [ ] Description uses retrieval language
- [ ] Examples show lookup and data usage patterns
- [ ] All lookup arguments documented
- [ ] All available attributes documented
- [ ] No import documentation included
- [ ] Only read timeout documented
- [ ] Practical lookup examples provided

This documentation should be treated as the source of truth for creating and maintaining high-quality, differentiated documentation for both Terraform AzureRM provider resources and data sources.

### Field Documentation Rules

#### Field Ordering Standards
- **Required fields first**: Always list required fields before optional fields in argument documentation
- **Alphabetical within category**: Within required and optional groups, list fields alphabetically
- **Consistent structure**: Maintain the same field ordering pattern across all block documentation

#### Note Format Standards
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

#### Azure-Specific Documentation Standards
- **Valid values only**: Only document values that are actually supported by the Azure service
- **API validation**: Verify all possible values against Azure SDK constants and API documentation
- **Cross-reference validation**: When implementing similar features across resources, ensure consistent value documentation
- **SDK alignment**: Match documentation values with Azure SDK enum constants where applicable

**When to Apply the "None" Pattern:**
- Optional fields that have Azure service defaults
- Fields where "None", "Off", or "Default" are valid Azure API values but add no user value
- Fields that follow the provider's move away from exposing Azure default constants

**When NOT to Apply the "None" Pattern:**
- Required fields
- Optional fields where the default behavior is not intuitive
- Fields where users commonly need to explicitly set values
