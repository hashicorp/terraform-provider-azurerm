---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_definition"
sidebar_current: "docs-azurerm-datasource-role-definition"
description: |-
  Get information about a custom Role Definition.
---

# Data Source: azurerm_role_definition

Use this data source to access the properties of a custom Role Definition. To access information about a built-in Role Definition, [please see the `azurerm_builtin_role_definition` data source](builtin_role_definition.html) instead.

## Example Usage

```hcl
data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "custom" {
  role_definition_id = "00000000-0000-0000-0000-000000000000"
  scope              = "${data.azurerm_subscription.primary.id}" # /subscriptions/00000000-0000-0000-0000-000000000000
}

output "custom_role_definition_id" {
  value = "${data.azurerm_role_definition.custom.id}"
}
```

## Argument Reference

* `role_definition_id` - (Required) Specifies the ID of the Role Definition as a UUID/GUID.

* `scope` - (Required) Specifies the Scope at which the Custom Role Definition exists.

## Attributes Reference

* `id` - the ID of the built-in Role Definition.
* `description` - the Description of the built-in Role.
* `type` - the Type of the Role.
* `permissions` - a `permissions` block as documented below.
* `assignable_scopes` - One or more assignable scopes for this Role Definition, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`.

A `permissions` block contains:

* `actions` - a list of actions supported by this role
* `not_actions` - a list of actions which are denied by this role
