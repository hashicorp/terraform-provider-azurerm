---
applyTo: "internal/website/docs/**/*.html.markdown"
description: This document outlines the standards and guidelines for writing documentation for Terraform resources and data sources in the AzureRM provider.
---

# Documentation Guidelines

This document outlines the standards and guidelines for writing documentation for Terraform resources and data sources in the AzureRM provider.

**Quick navigation:** [🚨 Pre-Implementation Requirements](#🚨-critical-pre-implementation-requirements-🚨) | [📚 Key Differences](#📚-key-differences-resources-vs-data-sources) | [🏗️ Documentation Structure](#🏗️-documentation-structure) | [📄 Resource Template](#📄-resource-documentation-template) | [📊 Data Source Template](#📊-data-source-documentation-template) | [✍️ Writing Guidelines](#✍️-writing-guidelines) | [💡 Example Configuration](#💡-example-configuration-guidelines) | [📁 Import Documentation](#📁-import-documentation) | [⏱️ Timeout Documentation](#⏱️-timeout-documentation) | [☁️ Azure-Specific Patterns](#☁️-azure-specific-documentation-patterns) | [📋 Attributes Reference](#📋-attributes-reference-differences) | [📝 Field Documentation](#📝-field-documentation-rules) | [📋 Provider Standards](#📋-provider-documentation-standards-note-formatting)

## 🚨 **CRITICAL: PRE-IMPLEMENTATION REQUIREMENTS** 🚨

**⚠️ MANDATORY BEFORE ANY DOCUMENTATION CHANGES ⚠️**

**BEFORE making ANY documentation changes, you MUST:**

1. **📋 READ NOTE FORMATTING GUIDELINES FIRST**
   - Scroll to [Provider Documentation Standards (Note Formatting)](#📋-provider-documentation-standards-note-formatting)
   - Review the three note types: Informational (`->`), Warning (`~>`), Caution (`!>`)
   - Understand the categorization criteria for each type

2. **🎯 CATEGORIZE YOUR CONTENT**
   - **Informational (`-> **Note:**`)**: Additional useful information, recommendations, tips, external links
   - **Warning (`~> **Note:**`)**: Information to avoid errors that won't cause irreversible changes (ForceNew behavior, conditional requirements)
   - **Caution (`!> **Note:**`)**: Critical information about irreversible changes, data loss, permanent effects

3. **✅ VALIDATE BEFORE IMPLEMENTATION**
   - Ask yourself: "What type of information am I documenting?"
   - Choose the appropriate note format based on the categorization criteria
   - ForceNew behavior = Warning note (`~> **Note:**`) - users need to avoid configuration errors
   - Azure service limitations = Often caution notes (`!> **Note:**`) if irreversible
   - Additional information/tips = Informational notes (`-> **Note:**`)

**🚫 COMMON MISTAKES TO AVOID:**
- Using informational notes (`->`) for ForceNew behavior warnings
- Using warning notes (`~>`) for simple tips or external links
- Using caution notes (`!>`) for reversible configuration changes

**📋 ENFORCEMENT CHECKLIST:**
- [ ] Read the note formatting guidelines section first
- [ ] Categorized the information type according to the criteria
- [ ] Chosen the appropriate note format based on impact and reversibility
- [ ] Verified the format matches the content type (warning for ForceNew, etc.)

---
[⬆️ Back to top](#documentation-guidelines)

## 📚 Key Differences: Resources vs Data Sources

**Language Patterns:**
- **Resources**: Use action verbs - `Manages`, `Creates`, `Configures`
- **Data Sources**: Use retrieval verbs - `Gets information about`, `Use this data source to access information about`

**Description Patterns:**
```markdown
# Resource Description
description: |-
  Manages a Service Resource.

# Data Source Description
description: |-
  Gets information about an existing Service Resource.
```

**Argument Types:**
- **Resources**: Arguments are for configuration (Required/Optional)
- **Data Sources**: Arguments are for identification/filtering (Required for lookup)

**Attributes Reference:**
- **Resources**: Exports computed values after creation/update
- **Data Sources**: Exports all available information from existing resources

**Timeout Blocks:**
- **Resources**: Include all CRUD operations
- **Data Sources**: Only include read operation

**Import Documentation:**
- **Resources**: Include import section with example
- **Data Sources**: Omit import section (data sources don't support import)

---
[⬆️ Back to top](#documentation-guidelines)

## 🏗️ Documentation Structure

**File Organization:**
```text
website/docs/
 r/                                # Resource documentation
    service_resource.html.markdown
 d/                                # Data source documentation
    service_resource.html.markdown
 guides/                           # Provider guides and tutorials
     guide_name.html.markdown
```

**File Naming:**
- **Resources**: `r/service_resourcetype.html.markdown`
- **Data Sources**: `d/service_resourcetype.html.markdown`
- Use lowercase with underscores, match Terraform resource name exactly

---
[⬆️ Back to top](#documentation-guidelines)

## 📄 Resource Documentation Template

### Standard Resource Documentation Structure

````markdown
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

* `enabled` - (Optional) Is this Service Resource enabled? Defaults to `true`.

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
````

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

---
[⬆️ Back to top](#documentation-guidelines)

## 📊 Data Source Documentation Template

### Standard Data Source Documentation Structure

````markdown
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
````

---
[⬆️ Back to top](#documentation-guidelines)

## ✍️ Writing Guidelines

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

### Nested Block Documentation
```markdown
---

A `configuration` block supports the following:

* `required_setting` - (Required) Description of the required setting.

* `optional_setting` - (Optional) Description of the optional setting. Defaults to `default_value`.
```

---
[⬆️ Back to top](#documentation-guidelines)

## 💡 Example Configuration Guidelines

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
- **Update existing**: `response_timeout_seconds = 120` (simple timeout field)
- **Update existing**: `enabled = false` (basic boolean toggle)
- **New example needed**: Complex nested `security_policy` block with multiple sub-configurations
- **New example needed**: Advanced `custom_domain` setup requiring certificates and DNS validation

---
[⬆️ Back to top](#documentation-guidelines)

## 📁 Import Documentation

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

---
[⬆️ Back to top](#documentation-guidelines)

## ⏱️ Timeout Documentation

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

---
[⬆️ Back to top](#documentation-guidelines)

## ☁️ Azure-Specific Documentation Patterns

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

---
[⬆️ Back to top](#documentation-guidelines)

## 📋 Attributes Reference Differences

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

* `enabled` - Is the Service Resource enabled?

* `configuration` - A `configuration` block as defined below.

* `endpoint` - The endpoint URL of the Service Resource.

* `tags` - A mapping of tags assigned to the resource.
```

---
[⬆️ Back to top](#documentation-guidelines)

## 📝 Field Documentation Rules

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

---
[⬆️ Back to top](#documentation-guidelines)

## 📋 Provider Documentation Standards (Note Formatting)

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

---
[⬆️ Back to top](#documentation-guidelines)

---

## Quick Reference Links

- 🏠 **Home**: [../copilot-instructions.md](../copilot-instructions.md)
- ☁️ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md)
- 📋 **Code Clarity Enforcement**: [code-clarity-enforcement.instructions.md](./code-clarity-enforcement.instructions.md)
- ❌ **Error Patterns**: [error-patterns.instructions.md](./error-patterns.instructions.md)
- 🏗️ **Implementation Guide**: [implementation-guide.instructions.md](./implementation-guide.instructions.md)
- 🔄 **Migration Guide**: [migration-guide.instructions.md](./migration-guide.instructions.md)
- 🏢 **Provider Guidelines**: [provider-guidelines.instructions.md](./provider-guidelines.instructions.md)
- 📐 **Schema Patterns**: [schema-patterns.instructions.md](./schema-patterns.instructions.md)
- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md)

### 🚀 Enhanced Guidance Files

- 🔄 **API Evolution**: [api-evolution-patterns.instructions.md](./api-evolution-patterns.instructions.md)
- ⚡ **Performance**: [performance-optimization.instructions.md](./performance-optimization.instructions.md)
- 🔐 **Security**: [security-compliance.instructions.md](./security-compliance.instructions.md)
- 🔧 **Troubleshooting**: [troubleshooting-decision-trees.instructions.md](./troubleshooting-decision-trees.instructions.md)

---
[⬆️ Back to top](#documentation-guidelines)
