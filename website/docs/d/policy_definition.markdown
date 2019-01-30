---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_definition"
sidebar_current: "docs-azurerm-datasource-policy-definition"
description: |-
  Get information about a Policy Definition.
---

# Data Source: azurerm_policy_definition

Use this data source to access information about a Policy Definition, both custom and built in.

## Example Usage

```hcl
data "azurerm_policy_definition" "test" {
  display_name = "Allowed resource types"
}

output "id" {
  value = "${data.azurerm_policy_definition.test.id}"
}
```

## Argument Reference

* `display_name` - (Required) Specifies the name of the Policy Definition.


## Attributes Reference

* `id` - the ID of the Policy Definition.
* `description` - the Description of the Policy.
* `type` - the Type of Policy.
* `name` - The Name of the Policy Definition
