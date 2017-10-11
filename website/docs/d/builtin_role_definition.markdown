---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_builtin_role_definition"
sidebar_current: "docs-azurerm-datasource-builtin-role-definition"
description: |-
  Get information about a built-in Role Definition.
---

# azurerm_built_in_role_definition

Use this data source to access the properties of a built-in Role Definition.

## Example Usage

```hcl
data "azurerm_builtin_role_definition" "contributor" {
  name = "Contributor"
}

output "contributor_role_definition_id" {
  value = "${data.azurerm_built_in_role.contributor.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the built-in Role Definition. Possible values are: `Contributor`, `Owner`, `Reader` and `VirtualMachineContributor`.


## Attributes Reference

* `id` - the ID of the built-in Role Definition.
* `description` - the Description of the built-in Role.
* `type` - the Type of the Role.
* `permissions` - a `permissions` block as documented below.

A `permissions` block contains:

* `actions` - a list of actions supported by this role
* `not_actions` - a list of actions which are denied by this role
