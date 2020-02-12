---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_definition"
description: |-
  Get information about a Policy Definition.
---

# Data Source: azurerm_policy_definition

Use this data source to access information about a Policy Definition, both custom and built in. Retrieves Policy Definitions from your current subscription by default.

## Example Usage

```hcl
data "azurerm_policy_definition" "example" {
  display_name = "Allowed resource types"
}

output "id" {
  value = data.azurerm_policy_definition.example.id
}
```

## Argument Reference

* `display_name` - Specifies the name of the Policy Definition.
* `management_group_id` - (Optional) Only retrieve Policy Definitions from this Management Group.


## Attributes Reference

* `id` - The ID of the Policy Definition.
* `name` - The Name of the Policy Definition.
* `type` - The Type of Policy.
* `description` - The Description of the Policy.
* `policy_type` - The Type of the Policy, such as `Microsoft.Authorization/policyDefinitions`.
* `policy_rule` - The Rule as defined (in JSON) in the Policy.
* `parameters` - Any Parameters defined in the Policy.
* `metadata` - Any Metadata defined in the Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Definition.
