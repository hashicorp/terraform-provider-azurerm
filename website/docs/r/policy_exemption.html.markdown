---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_exemption"
description: |-
  Manages a policy exemption
---

# azurerm_policy_exemption

Exempts a resource hierarchy or an individual resource from evaluation of initiatives or definitions.

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

resource "azurerm_policy_exemption" "test" {
  name                 = "example-policy-exemption"
  scope                = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id

  exemption_category = "Mitigated"

  display_name = "My example policy exemption"
  description  = "Policy Exemption created as an example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Policy Exemption. Changing this forces a new resource to be created.

* `scope`- (Required) The Scope at which the Policy Exemption should be applied, which must be a Resource ID (such as Subscription e.g. `/subscriptions/00000000-0000-0000-000000000000` or a Resource Group e.g.`/subscriptions/00000000-0000-0000-000000000000/resourceGroups/myResourceGroup`). Changing this forces a new resource to be created.

* `exemption_category` - (Required) The category of this policy exemption. Possible values are `Waiver` and `Mitigated`.

* `policy_assignment_id` - (Required) The ID of the Policy Assignment to be exempted at the specified Scope.

* `description` - (Optional) A description to use for this Policy Exemption.

* `display_name` - (Optional) A friendly display name to use for this Policy Exemption.

* `expires_on` - (Optional) The expiration date and time in UTC ISO 8601 format of this policy exemption.

* `policy_definition_reference_ids` - (Optional) The policy definition reference ID list when the associated policy assignment is an assignment of a policy set definition.

* `metadata` - (Optional) The metadata for this policy exemption. This is a JSON string representing additional metadata that should be stored with the policy exemption.

## Attributes Reference

The following attributes are exported:

* `id` - The Policy Exemption id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Exemption.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Exemption.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Exemption.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Exemption.

## Import

Policy Exemptions can be imported using the `policy name`, e.g.

```shell
terraform import azurerm_policy_exemption.assignment1  /subscriptions/00000000-0000-0000-000000000000/providers/Microsoft.Authorization/policyExemptions/exemption1
```
