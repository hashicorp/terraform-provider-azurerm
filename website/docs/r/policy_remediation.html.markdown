---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_remediation"
description: |-
  Manages an Azure Policy Remediation.
---

# azurerm_policy_insights_remediation

Manages an Azure Policy Remediation at the specified Scope.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_policy_definition" "example" {
  name         = "my-policy-definition"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "my-policy-definition"

  policy_rule = <<POLICY_RULE
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

resource "azurerm_policy_assignment" "example" {
  name                 = "example-policy-assignment"
  scope                = azurerm_resource_group.example.id
  policy_definition_id = azurerm_policy_definition.example.id
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "My Example Policy Assignment"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West Europe" ]
  }
}
PARAMETERS
}

resource "azurerm_policy_remediation" "example" {
  name                 = "example-policy-remediation"
  scope                = azurerm_policy_assignment.example.scope
  policy_assignment_id = azurerm_policy_assignment.example.id
  location_filters     = ["West Europe"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Policy Remediation. Changing this forces a new resource to be created.

* `scope` - (Required) The Scope at which the Policy Remediation should be applied. Changing this forces a new resource to be created. A scope must be a Resource ID out of one of the following list:

    1. A subscription, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000`
    1. A resource group, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup`
    1. A resource, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM`
    1. A management group, e.g. `/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000`

* `policy_assignment_id` - (Required) The ID of the Policy Assignment that should be remediated.

* `policy_definition_reference_id` - (Optional) The unique ID for the policy definition within the policy set definition that should be remediated. Required when the policy assignment being remediated assigns a policy set definition.

* `location_filters` - (Optional) A list of the resource locations that will be remediated.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Policy Remediation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Remediation.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Remediation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Remediation.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Remediation.


## Import

Policy Remediations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_policy_remediation.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.PolicyInsights/remediations/remediation1
```
or
```shell
terraform import azurerm_policy_remediation.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.PolicyInsights/remediations/remediation1
```
or
```shell
terraform import azurerm_policy_remediation.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.PolicyInsights/remediations/remediation1
```
or
```shell
terraform import azurerm_policy_remediation.example /providers/Microsoft.Management/managementGroups/my-mgmt-group-id/providers/Microsoft.PolicyInsights/remediations/remediation1
```
