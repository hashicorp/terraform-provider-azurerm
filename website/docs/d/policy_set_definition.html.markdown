---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_set_definition"
description: |-
  Get information about a Policy Set Definition.
---

# Data Source: azurerm_policy_set_definition

Use this data source to access information about an existing Policy Set Definition.

## Example Usage

```hcl
data "azurerm_policy_set_definition" "example" {
  display_name = "Policy Set Definition Example"
}

output "id" {
  value = data.azurerm_policy_set_definition.example.id
}
```

## Argument Reference

* `name` - Specifies the name of the Policy Set Definition. Conflicts with `display_name`.

* `display_name` - Specifies the display name of the Policy Set Definition. Conflicts with `name`.

**NOTE** As `display_name` is not unique errors may occur when there are multiple policy set definitions with same display name. 

* `management_group_name` - (Optional) Only retrieve Policy Set Definitions from this Management Group.

## Attributes Reference

* `id` - The ID of the Policy Set Definition.

* `description` - The Description of the Policy Set Definition.

* `policy_type` - The Type of the Policy Set Definition.

* `policy_definitions` - The policy definitions contained within the policy set definition.

* `policy_definition_reference` - One or more `policy_definition_reference` blocks as defined below.

* `policy_definition_group` - One or more `policy_definition_group` blocks as defined below.

* `parameters` - Any Parameters defined in the Policy Set Definition.

* `metadata` - Any Metadata defined in the Policy Set Definition.

---

An `policy_definition_reference` block exports the following:

* `policy_definition_id` - The ID of the policy definition or policy set definition that is included in this policy set definition.

* `parameters` - The mapping of the parameter values for the referenced policy rule. The keys are the parameter names.

-> **NOTE:** Since Terraform's concept of a map requires all of the elements to be of the same type, the value in parameters will all be converted to string type.

~> **Note:** This field only supports String fields and is deprecated in favour of the `parameters_values` field

* `parameter_values` - The parameter values for the referenced policy rule. This field is a json object.

* `reference_id` - The unique ID within this policy set definition for this policy definition reference.

* `group_names` - The list of names of the policy definition groups that this policy definition reference belongs to.

---

An `policy_definition_group` block exports the following:

* `name` - The name of this policy definition group.

* `display_name` - The display name of this policy definition group. 

* `category` - The category of this policy definition group.

* `description` - The description of this policy definition group.

* `additional_metadata_id` - The ID of a resource that contains additional metadata about this policy definition group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Set Definition.
