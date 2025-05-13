---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_group_policy_assignment"
description: |-
  Manages a Resource Group Policy Assignment.
---

# azurerm_resource_group_policy_assignment

Manages a Resource Group Policy Assignment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_policy_definition" "example" {
  name         = "only-deploy-in-westeurope"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "my-policy-definition"

  policy_rule = <<POLICY_RULE
 {
    "if": {
      "not": {
        "field": "location",
        "equals": "westeurope"
      }
    },
    "then": {
      "effect": "Deny"
    }
  }
POLICY_RULE
}

resource "azurerm_resource_group_policy_assignment" "example" {
  name                 = "example"
  resource_group_id    = azurerm_resource_group.example.id
  policy_definition_id = azurerm_policy_definition.example.id

  parameters = <<PARAMS
    {
      "tagName": {
        "value": "Business Unit"
      },
      "tagValue": {
        "value": "BU"
      }
    }
PARAMS
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Policy Assignment. Changing this forces a new Policy Assignment to be created. Cannot exceed 64 characters in length.

* `policy_definition_id` - (Required) The ID of the Policy Definition or Policy Definition Set. Changing this forces a new Policy Assignment to be created.

* `resource_group_id` - (Required) The ID of the Resource Group where this Policy Assignment should be created. Changing this forces a new Policy Assignment to be created.

---

* `description` - (Optional) A description which should be used for this Policy Assignment.

* `display_name` - (Optional) The Display Name for this Policy Assignment.

* `enforce` - (Optional) Specifies if this Policy should be enforced or not? Defaults to `true`.

* `identity` - (Optional) An `identity` block as defined below.

-> **Note:** The `location` field must also be specified when `identity` is specified.

* `location` - (Optional) The Azure Region where the Policy Assignment should exist. Changing this forces a new Policy Assignment to be created.

* `metadata` - (Optional) A JSON mapping of any Metadata for this Policy.

* `non_compliance_message` - (Optional) One or more `non_compliance_message` blocks as defined below.

* `not_scopes` - (Optional) Specifies a list of Resource Scopes (for example a Subscription, or a Resource Group) within this Management Group which are excluded from this Policy.

* `parameters` - (Optional) A JSON mapping of any Parameters for this Policy.

* `overrides` - (Optional) One or more `overrides` blocks as defined below. More detail about `overrides` and `resource_selectors` see [policy assignment structure](https://learn.microsoft.com/en-us/azure/governance/policy/concepts/assignment-structure#resource-selectors-preview)

* `resource_selectors` - (Optional) One or more `resource_selectors` blocks as defined below to filter polices by resource properties.

---

A `identity` block supports the following:

* `type` - (Required) The Type of Managed Identity which should be added to this Policy Definition. Possible values are `SystemAssigned` and `UserAssigned`.

* `identity_ids` - (Optional) A list of User Managed Identity IDs which should be assigned to the Policy Definition.

~> **Note:** This is required when `type` is set to `UserAssigned`.

---

A `non_compliance_message` block supports the following:

* `content` - (Required) The non-compliance message text. When assigning policy sets (initiatives), unless `policy_definition_reference_id` is specified then this message will be the default for all policies.

* `policy_definition_reference_id` - (Optional) When assigning policy sets (initiatives), this is the ID of the policy definition that the non-compliance message applies to.

---

A `overrides` block supports the following:

* `value` - (Required) Specifies the value to override the policy property. Possible values for `policyEffect` override listed [policy effects](https://learn.microsoft.com/en-us/azure/governance/policy/concepts/effects).

* `selectors` - (Optional) One or more `override_selector` block as defined below.

---

A `override_selector` block supports the following:

* `in` - (Optional) Specify the list of policy reference id values to filter in. Cannot be used with `not_in`.

* `not_in` - (Optional) Specify the list of policy reference id values to filter out. Cannot be used with `in`.

---

A `resource_selectors` block supports the following:

* `name` - (Optional) Specifies a name for the resource selector.

* `selectors` - (Required) One or more `resource_selector` block as defined below.

---

A `resource_selector` block supports the following:

* `kind` - (Required) Specifies which characteristic will narrow down the set of evaluated resources. Possible values are `resourceLocation`, `resourceType` and `resourceWithoutLocation`.

* `in` - (Optional) The list of allowed values for the specified kind. Cannot be used with `not_in`. Can contain up to 50 values.

* `not_in` - (Optional) The list of not-allowed values for the specified kind. Cannot be used with `in`. Can contain up to 50 values.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Group Policy Assignment.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID of the Policy Assignment for this Resource Group.

* `tenant_id` - The Tenant ID of the Policy Assignment for this Resource Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Assignment for this Resource Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Assignment for this Resource Group.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Assignment for this Resource Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Assignment for this Resource Group.

## Import

Resource Group Policy Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_group_policy_assignment.example /subscriptions/00000000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Authorization/policyAssignments/assignment1
```
