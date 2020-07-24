---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_set_definition"
description: |-
  Get information about a Policy Set Definition.
---

# Data Source: azurerm_policy_set_definition

Use this data source to access information about an existing Policy Set Definition.

## Example Usage

```hcl
data "azurerm_policy_set_definition" "example" {
  display_name = "Policy Set Definition Example"
}

output "id" {
  value = data.azurerm_policy_set_definition.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Policy Set Definition. Conflicts with `display_name`.

* `display_name` - Specifies the display name of the Policy Set Definition. Conflicts with `name`.

**NOTE** As `display_name` is not unique errors may occur when there are multiple policy set definitions with same display name. 

* `management_group_name` - (Optional) Only retrieve Policy Set Definitions from this Management Group.

## Attributes Reference

* `id` - The ID of the Policy Set Definition.

* `description` - The Description of the Policy Set Definition.

* `policy_type` - The Type of the Policy Set Definition.

* `policy_definitions` - The policy definitions contained within the policy set definition.

* `parameters` - Any Parameters defined in the Policy Set Definition.

* `metadata` - Any Metadata defined in the Policy Set Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Set Definition.
