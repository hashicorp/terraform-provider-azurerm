---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group_policy_set_definition"
description: |-
  Manages a Policy Set Definition for a Management Group.
---

# azurerm_management_group_policy_set_definition

Manages a Policy Set Definition for a Management Group.

## Example Usage

```hcl
resource "azurerm_management_group" "example" {
  display_name = "Example"
}

resource "azurerm_management_group_policy_set_definition" "example" {
  name                = "example"
  policy_type         = "Custom"
  display_name        = "Example"
  management_group_id = azurerm_management_group.example.id

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
    parameter_values     = <<VALUES
   {
     "listOfAllowedLocations": {"value": "[parameters('allowedLocations')]"}
   }
VALUES
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Policy Set Definition. Changing this forces a new Policy Set Definition to be created.
 
* `management_group_id` - (Required) The ID of the Management Group where this Policy Set Definition should be created. Changing this forces a new Policy Set Definition to be created.

* `display_name` - (Required) The display name of this Policy Set Definition.

* `policy_definition_reference` - (Required) One or more `policy_definition_reference` blocks as defined below.

* `policy_type` - (Required) The Policy Set Definition type. Possible values are `BuiltIn`, `Custom`, `NotSpecified`, and `Static`. Changing this forces a new Policy Set Definition to be created.

---

* `description` - (Optional) The description of this Policy Set Definition.

* `metadata` - (Optional) The metadata for the Policy Set Definition in JSON format.

* `parameters` - (Optional) The parameters for the Policy Set Definition in JSON format.

* `policy_definition_group` - (Optional) One or more `policy_definition_group` blocks as defined below.

---

A `policy_definition_group` block supports the following:

* `name` - (Required) The name which should be used for this Policy Definition Group.

* `additional_metadata_resource_id` - (Optional) The ID of a resource that contains additional metadata for this Policy Definition Group.

* `category` - (Optional) The category of this Policy Definition Group.

* `description` - (Optional) The description of this Policy Definition Group.

* `display_name` - (Optional) The display name of this Policy Definition Group.

---

A `policy_definition_reference` block supports the following:

* `policy_definition_id` - (Required) The ID of the Policy Definition to include in this Policy Set Definition.

* `parameter_values` - (Optional) Parameter values for the references Policy Definition in JSON format.

* `policy_group_names` - (Optional) Specifies a list of Policy Definition Groups names that this Policy Definition Reference belongs to.

* `reference_id` - (Optional) A unique ID within this Policy Set Definition for this Policy Definition Reference.

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
terraform import azurerm_management_group_policy_set_definition.example /providers/Microsoft.Management/managementGroups/0000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policySetDefinitions/policySetDefinitionName
```
