---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_definition"
**sidebar_current**: "docs-azurerm-datasource-policy-definition"
description: |-
  Get information about a Policy Definition.
---

# Data Source: azurerm_policy_definition

Use this data source to access information about a Policy Definition, both custom and built in. Looks at your current subscription by default.

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
* `management_group_id` - (Optional) look at this management_group for the policy definitions.


## Attributes Reference

* `id` - the ID of the Policy Definition.
* `name` - The Name of the Policy Definition.
* `type` - the Type of Policy.
* `description` - the Description of the Policy.
* `policy_type` - the Type of the Policy, such as `Microsoft.Authorization/policyDefinitions`.
* `policy_rule` - the Rule as defined (in JSON) in the Policy.
* `parameters` - any Parameters defined in the Policy.
* `metadata` - any Metadata defined in the Policy.
