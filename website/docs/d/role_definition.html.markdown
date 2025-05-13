---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_role_definition"
description: |-
  Get information about an existing Role Definition.
---

# Data Source: azurerm_role_definition

Use this data source to access information about an existing Role Definition.

## Example Usage

```hcl
data "azurerm_subscription" "primary" {
}

resource "azurerm_role_definition" "custom" {
  role_definition_id = "00000000-0000-0000-0000-000000000000"
  name               = "CustomRoleDef"
  scope              = data.azurerm_subscription.primary.id
  #...
}

data "azurerm_role_definition" "custom" {
  role_definition_id = azurerm_role_definition.custom.role_definition_id
  scope              = data.azurerm_subscription.primary.id # /subscriptions/00000000-0000-0000-0000-000000000000
}

data "azurerm_role_definition" "custom-byname" {
  name  = azurerm_role_definition.custom.name
  scope = data.azurerm_subscription.primary.id
}

data "azurerm_role_definition" "builtin" {
  name = "Contributor"
}

output "custom_role_definition_id" {
  value = data.azurerm_role_definition.custom.id
}

output "contributor_role_definition_id" {
  value = data.azurerm_role_definition.builtin.id
}
```

## Argument Reference

* `name` - (Optional) Specifies the Name of either a built-in or custom Role Definition.

-> **Note:** You can also use this for built-in roles such as `Contributor`, `Owner`, `Reader` and `Virtual Machine Contributor`

* `role_definition_id` - (Optional) Specifies the ID of the Role Definition as a UUID/GUID.

* `scope` - (Optional) Specifies the Scope at which the Custom Role Definition exists.

~> **Note:** One of `name` or `role_definition_id` must be specified.

## Attributes Reference

* `id` - The ID of the built-in Role Definition.

* `description` - The Description of the built-in Role.

* `type` - The Type of the Role.

* `permissions` - A `permissions` block as documented below.

* `assignable_scopes` - One or more assignable scopes for this Role Definition, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`.

---

A `permissions` block contains:

* `actions` - A list of actions supported by this role.

* `not_actions` - A list of actions which are denied by this role.

* `data_actions` - A list of data actions allowed by this role.

* `not_data_actions` - A list of data actions which are denied by this role.

* `condition` - The conditions on this role definition, which limits the resources it can be assigned to.

* `condition_version` - The version of the condition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Role Definition.
