---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_definition"
description: |-
  Manages a policy rule definition. Policy definitions do not take effect until they are assigned to a scope using a Policy Assignment.
---

# azurerm_policy_definition

Manages a policy rule definition on a management group or your provider subscription.

Policy definitions do not take effect until they are assigned to a scope using a Policy Assignment.

## Example Usage

```hcl
resource "azurerm_policy_definition" "policy" {
  name         = "accTestPolicy"
  policy_type  = "Custom"
  mode         = "Indexed"
  display_name = "acceptance test policy definition"

  metadata = <<METADATA
    {
    "category": "General"
    }

METADATA


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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy definition. Changing this forces a
    new resource to be created.

* `policy_type` - (Required) The policy type. Possible values are `BuiltIn`, `Custom` and `NotSpecified`. Changing this forces a new resource to be created.

* `mode` - (Required) The policy mode that allows you to specify which resource
    types will be evaluated. Possible values are `All`, `Indexed`, `Microsoft.ContainerService.Data`, `Microsoft.CustomerLockbox.Data`, `Microsoft.DataCatalog.Data`, `Microsoft.KeyVault.Data`, `Microsoft.Kubernetes.Data`, `Microsoft.MachineLearningServices.Data`, `Microsoft.Network.Data` and `Microsoft.Synapse.Data`.

* `display_name` - (Required) The display name of the policy definition.

* `description` - (Optional) The description of the policy definition.

* `management_group_name` - (Optional) The name of the Management Group where this policy should be defined. Changing this forces a new resource to be created.

* `management_group_id` - (Optional / **Deprecated in favour of `management_group_name`**) The name of the Management Group where this policy should be defined. Changing this forces a new resource to be created.

~> **Note:** if you are using `azurerm_management_group` to assign a value to `management_group_id`, be sure to use `name` or `group_id` attribute, but not `id`.

* `policy_rule` - (Optional) The policy rule for the policy definition. This
    is a JSON string representing the rule that contains an if and
    a then block.

* `metadata` - (Optional) The metadata for the policy definition. This
    is a JSON string representing additional metadata that should be stored
    with the policy definition.

* `parameters` - (Optional) Parameters for the policy definition. This field
    is a JSON string that allows you to parameterize your policy definition.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Policy Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Definition.

## Import

Policy Definitions can be imported using the `policy name`, e.g.

```shell
terraform import azurerm_policy_definition.examplePolicy /subscriptions/<SUBSCRIPTION_ID>/providers/Microsoft.Authorization/policyDefinitions/<POLICY_NAME>
```

or

```shell
terraform import azurerm_policy_definition.examplePolicy /providers/Microsoft.Management/managementgroups/<MANGAGEMENT_GROUP_ID>/providers/Microsoft.Authorization/policyDefinitions/<POLICY_NAME>
```
