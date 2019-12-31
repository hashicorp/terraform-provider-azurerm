---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_remediation"
sidebar_current: "docs-azurerm-resource-policy-remediation"
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
  name = "example-policy-remediation"
  scope = azurerm_policy_assignment.example.scope
  policy_assignment_id = azurerm_policy_assignment.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Policy Remediation. Changing this forces a new resource to be created.

* `scope` - (Required) The Scope at which the Policy Remediation should be applied, which must be a Resource ID (such as Subscription e.g. `/subscriptions/00000000-0000-0000-0000-000000000000` or a Resource Group e.g.`/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup` or a specified Resource e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM`) or a Management Group ID (e.g. `/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000`). Changing this forces a new resource to be created.

* `policy_assignment_id` - (Required) The resource ID of the policy assignment that should be remediated.

* `policy_definition_reference_id` - (Optional) The policy definition reference ID of the individual definition that should be remediated. Required when the policy assignment being remediated assigns a policy set definition.

* `location_filters` - (Optional) A list of the resource locations that will be remediated.

## Attributes Reference

The following attributes are exported:

* `created_on` - The time at which the remediation was created.

* `last_updated_on` - The time at which the remediation was last updated.

* `id` - The ID of the Policy Remediation.

## Import

Policy Remediations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_policy_remediation.remediation1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.PolicyInsights/remediations/remediation1
```
