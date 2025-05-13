---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_set_definition"
description: |-
  Manages a policy set definition.
---

# azurerm_policy_set_definition

Manages a policy set definition.

-> **Note:** Policy set definitions (also known as policy initiatives) do not take effect until they are assigned to a scope using a Policy Set Assignment.

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

  policy_definition_reference {
    policy_definition_id = "/providers/Microsoft.Authorization/policyDefinitions/e765b5de-1225-4ba3-bd56-1ac6695af988"
    parameter_values     = <<VALUE
    {
      "listOfAllowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
    VALUE
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy set definition. Changing this forces a new resource to be created.

* `policy_type` - (Required) The policy set type. Possible values are `BuiltIn`, `Custom`, `NotSpecified` and `Static`. Changing this forces a new resource to be created.

* `display_name` - (Required) The display name of the policy set definition.

* `policy_definition_reference` - (Required) One or more `policy_definition_reference` blocks as defined below.

* `policy_definition_group` - (Optional) One or more `policy_definition_group` blocks as defined below.

* `description` - (Optional) The description of the policy set definition.

* `management_group_id` - (Optional) The id of the Management Group where this policy set definition should be defined. Changing this forces a new resource to be created.

* `metadata` - (Optional) The metadata for the policy set definition. This is a JSON object representing additional metadata that should be stored with the policy definition.

* `parameters` - (Optional) Parameters for the policy set definition. This field is a JSON object that allows you to parameterize your policy definition.

---

A `policy_definition_reference` block supports the following:

* `policy_definition_id` - (Required) The ID of the policy definition that will be included in this policy set definition.

* `parameter_values` - (Optional) Parameter values for the referenced policy rule. This field is a JSON string that allows you to assign parameters to this policy rule.

* `reference_id` - (Optional) A unique ID within this policy set definition for this policy definition reference.

* `policy_group_names` - (Optional) A list of names of the policy definition groups that this policy definition reference belongs to.

---

An `policy_definition_group` block supports the following:

* `name` - (Required) The name of this policy definition group.

* `display_name` - (Optional) The display name of this policy definition group.

* `category` - (Optional) The category of this policy definition group.

* `description` - (Optional) The description of this policy definition group.

* `additional_metadata_resource_id` - (Optional) The ID of a resource that contains additional metadata about this policy definition group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Policy Set Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Set Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Set Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Set Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Set Definition.

## Import

Policy Set Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_policy_set_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policySetDefinitions/testPolicySet
```

or

```shell
terraform import azurerm_policy_set_definition.example /providers/Microsoft.Management/managementGroups/my-mgmt-group-id/providers/Microsoft.Authorization/policySetDefinitions/testPolicySet
```
