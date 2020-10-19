---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_assignment"
description: |-
  Configures the specified Policy Definition at the specified Scope.
---

# azurerm_policy_assignment

Configures the specified Policy Definition at the specified Scope. Also, Policy Set Definitions are supported.

## Example Usage

```hcl
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

resource "azurerm_resource_group" "example" {
  name     = "test-resources"
  location = "West Europe"
}

resource "azurerm_policy_assignment" "example" {
  name                 = "example-policy-assignment"
  scope                = azurerm_resource_group.example.id
  policy_definition_id = azurerm_policy_definition.example.id
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "My Example Policy Assignment"

  metadata = <<METADATA
    {
    "category": "General"
    }
METADATA

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

* `scope`- (Required) The Scope at which the Policy Assignment should be applied, which must be a Resource ID (such as Subscription e.g. `/subscriptions/00000000-0000-0000-000000000000` or a Resource Group e.g.`/subscriptions/00000000-0000-0000-000000000000/resourceGroups/myResourceGroup`). Changing this forces a new resource to be created.

* `policy_definition_id` - (Required) The ID of the Policy Definition to be applied at the specified Scope.

* `identity` - (Optional) An `identity` block.

* `location` - (Optional) The Azure location where this policy assignment should exist. This is required when an Identity is assigned. Changing this forces a new resource to be created.

* `description` - (Optional) A description to use for this Policy Assignment. Changing this forces a new resource to be created.

* `display_name` - (Optional) A friendly display name to use for this Policy Assignment. Changing this forces a new resource to be created.

* `metadata` - (Optional) The metadata for the policy assignment. This is a json object representing additional metadata that should be stored with the policy assignment.

* `parameters` - (Optional) Parameters for the policy definition. This field is a JSON object that maps to the Parameters field from the Policy Definition. Changing this forces a new resource to be created.

~> **NOTE:** This value is required when the specified Policy Definition contains the `parameters` field.

* `not_scopes` - (Optional) A list of the Policy Assignment's excluded scopes. The list must contain Resource IDs (such as Subscriptions e.g. `/subscriptions/00000000-0000-0000-000000000000` or Resource Groups e.g.`/subscriptions/00000000-0000-0000-000000000000/resourceGroups/myResourceGroup`).

* `enforcement_mode`- (Optional) Can be set to 'true' or 'false' to control whether the assignment is enforced (true) or not (false). Default is 'true'.

---

An `identity` block supports the following:

* `type` - (Required) The Managed Service Identity Type of this Policy Assignment. Possible values are `SystemAssigned` (where Azure will generate a Service Principal for you), or `None` (no use of a Managed Service Identity).

~> **NOTE:** When `type` is set to `SystemAssigned`, identity the Principal ID can be retrieved after the policy has been assigned.

---


## Attributes Reference

The following attributes are exported:

* `id` - The Policy Assignment id.

* `identity` - An `identity` block.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID of this Policy Assignment if `type` is `SystemAssigned`.

* `tenant_id` - The Tenant ID of this Policy Assignment if `type` is `SystemAssigned`.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Assignment.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Assignment.

## Import

Policy Assignments can be imported using the `policy name`, e.g.

```shell
terraform import azurerm_policy_assignment.assignment1  /subscriptions/00000000-0000-0000-000000000000/providers/Microsoft.Authorization/policyAssignments/assignment1
```
