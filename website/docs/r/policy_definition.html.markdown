---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_definition"
sidebar_current: "docs-azurerm-resource-policy-definition"
description: |-
  Manages a policy for all of the resource groups under the subscription.
---

# azurerm_policy_definition

Manages a policy for all of the resource groups under the subscription.

## Example Usage

```hcl
resource "azurerm_policy_definition" "policy" {
  name         = "accTestPolicy"
  policy_type  = "BuiltIn"
  mode         = "Indexed"
  display_name = "acceptance test policy definition"
  policy_rule  = <<POLICY_RULE
	{
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
	{
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy definition. Changing this forces a
    new resource to be created.

* `policy_type` - (Required) The policy type.  The value can be "BuiltIn", "Custom"
    or "NotSpecified". Changing this forces a new resource to be created.

* `mode` - (Required) The policy mode that allows you to specify which resource
    types will be evaluated.  The value can be "All", "Indexed" or
    "NotSpecified". Changing this resource forces a new resource to be
    created.

* `display_name` - (Required) The display name of the policy definition.

* `description` - (Optional) The description of the policy definition.

* `policy_rule` - (Optional) The policy rule for the policy definition. This
    is a json object representing the rule that contains an if and
    a then block.

* `metadata` - (Optional) The metadata for the policy definition. This
    is a json object representing the rule that contains an if and
    a then block.

* `parameters` - (Optional) Parameters for the policy definition. This field
    is a json object that allows you to parameterize your policy definition.

## Attributes Reference

The following attributes are exported:

* `id` - The policy definition id.

## Import

Policy Definitions can be imported using the `policy name`, e.g.

```shell
terraform import azurerm_policy_definition.testPolicy  /subscriptions/<SUBSCRIPTION_ID>/providers/Microsoft.Authorization/policyDefinitions/<POLICY_NAME>
```
