---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_set_definition"
description: |-
  Manages a policy set definition.
---

# azurerm_policy_set_definition

Manages a policy set definition.

-> **NOTE:**  Policy set definitions (also known as policy initiatives) do not take effect until they are assigned to a scope using a Policy Set Assignment.

## Example Usage

```hcl
resource "azurerm_policy_set_definition" "example" {
  name         = "testPolicySet"
  policy_type  = "Custom"
  display_name = "Test Policy Set"

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


  policy_definitions = <<POLICY_DEFINITIONS
    [
        {
            "parameters": {
                "listOfAllowedLocations": {
                    "value": "[parameters('allowedLocations')]"
                }
            },
            "policyDefinitionId": "/providers/Microsoft.Authorization/policyDefinitions/e765b5de-1225-4ba3-bd56-1ac6695af988"
        }
    ]
POLICY_DEFINITIONS

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy set definition. Changing this forces a new resource to be created.

* `policy_type` - (Required) The policy set type. Possible values are `BuiltIn` or `Custom`. Changing this forces a new resource to be created.

* `display_name` - (Required) The display name of the policy set definition.

* `policy_definitions` - (Required) The policy definitions for the policy set definition. This is a json object representing the bundled policy definitions .

* `description` - (Optional) The description of the policy set definition.

* `management_group_id` - (Optional) The ID of the Management Group where this policy should be defined. Changing this forces a new resource to be created.

~> **Note:** if you are using `azurerm_management_group` to assign a value to `management_group_id`, be sure to use `.group_id` and not `.id`.

* `metadata` - (Optional) The metadata for the policy set definition. This is a json object representing additional metadata that should be stored with the policy definition.

* `parameters` - (Optional) Parameters for the policy set definition. This field is a json object that allows you to parameterize your policy definition.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Policy Set Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Set Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Set Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Set Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Set Definition.

## Import

Policy Set Definitions can be imported using the Resource ID, e.g.

```shell
terraform import azurerm_policy_set_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policySetDefinitions/testPolicySet
```
or
```shell
terraform import azurerm_policy_set_definition.example/providers/Microsoft.Management/managementgroups/my-mgmt-group-id/providers/Microsoft.Authorization/policySetDefinitions/testPolicySet
```
