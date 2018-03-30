---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_assignment"
sidebar_current: "docs-azurerm-resource-policy-assignment"
description: |-
  Configured the specified Policy Definition at the specified Scope.
---

# azurerm_policy_assignment

Configured the specified Policy Definition at the specified Scope.

## Example Usage

```hcl
resource "azurerm_policy_definition" "test" {
  name         = "my-policy-definition"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%d"
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

resource "azurerm_resource_group" "test" {
  name = "test-resources"
  location = "West Europe"
}

resource "azurerm_policy_assignment" "test" {
  name                 = "example-policy-assignment"
  scope                = "${azurerm_resource_group.test.id}"
  policy_definition_id = "${azurerm_policy_definition.test.id}"
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "Acceptance Test Run %d"
  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West Europe" ]
  }
}
PARAMETERS
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Policy Assignment. Changing this forces a new resource to be created.

* `scope`- (Required) The Scope at which the Policy Assignment should be applied. This can either be the Subscription (e.g. `/subscriptions/00000000-0000-0000-000000000000`) or a Resource Group (e.g.`/subscriptions/00000000-0000-0000-000000000000/resourceGroups/myResourceGroup`). Changing this forces a new resource to be created.

* `policy_definition_id` - (Required) The ID of the Policy Definition to be applied at the specified Scope.

* `description` - (Optional) A description to use for this Policy Assignment. Changing this forces a new resource to be created.

* `display_name` - (Optional) A friendly display name to use for this Policy Assignment. Changing this forces a new resource to be created.

* `parameters` - (Optional) Parameters for the policy definition. This field is a JSON object that maps to the Parameters field from the Policy Definition. Changing this forces a new resource to be created.

~> **NOTE:** This value is required when the specified Policy Definition contains the `parameters` field.

## Attributes Reference

The following attributes are exported:

* `id` - The Policy Assignment id.

## Import

Policy Assignments can be imported using the `policy name`, e.g.

```shell
terraform import azurerm_policy_assignment.assignment1  /subscriptions/00000000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/assignment1
```
