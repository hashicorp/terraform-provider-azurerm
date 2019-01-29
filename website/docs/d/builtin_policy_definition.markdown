---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_builtin_policy_definition"
sidebar_current: "docs-azurerm-datasource-builtin-policy-definition"
description: |-
  Get information about an existing built-in Policy Definition.
---

# Data Source: azurerm_builtin_policy_definition

Use this data source to access information about a built-in Policy Definition. To access information about a custom Policy Definition, [please see the `azurerm_policy_definition` data source](../r/policy_definition.html) instead.

## Example Usage

```hcl
data "azurerm_builtin_policy_definition" "test" {
  display_name = "Allowed resource types"
}

output "id" {
  value = "${data.azurerm_builtin_policy_definition.test.id}"
}
```

## Argument Reference

* `display_name` - (Required) Specifies the name of the built-in Policy Definition.


## Attributes Reference

* `id` - the ID of the built-in Policy Definition.
* `description` - the Description of the built-in Policy.
* `type` - the Type of Policy.
* `name` - The Name of the Policy Definition