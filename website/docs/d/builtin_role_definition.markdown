---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_builtin_role_definition"
sidebar_current: "docs-azurerm-datasource-builtin-role-definition"
description: |-
  Get information about an existing built-in Role Definition.
---

# Data Source: azurerm_builtin_role_definition

Use this data source to access information about a built-in Role Definition. To access information about a custom Role Definition, [please see the `azurerm_role_definition` data source](role_definition.html) instead.

~> **NOTE:** The this datasource has been deprecated in favour of `azurerm_role_definition` that now can look up role definitions by name. As such this data source will be removed in version 2.0 of the AzureRM Provider.

## Example Usage

```hcl
data "azurerm_builtin_role_definition" "contributor" {
  name = "Contributor"
}

output "contributor_role_definition_id" {
  value = "${data.azurerm_builtin_role_definition.contributor.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the built-in Role Definition. Possible values are: `Contributor`, `Owner`, `Reader` and `VirtualMachineContributor`.


## Attributes Reference

* `id` - the ID of the built-in Role Definition.
* `description` - the Description of the built-in Role.
* `type` - the Type of the Role.
* `permissions` - a `permissions` block as documented below.
* `assignable_scopes` - One or more assignable scopes for this Role Definition, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`.

A `permissions` block contains:

* `actions` - a list of actions supported by this role
* `data_actions` - a list of data actions supported by this role
* `not_actions` - a list of actions which are denied by this role
* `not_data_actions` - a list of data actions which are denied by this role
