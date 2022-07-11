---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group_policy_assignment"
description: |-
  Manages a Policy Assignment to a Management Group.
---

# azurerm_management_group_policy_assignment

Manages a Policy Assignment to a Management Group.

## Example Usage

```hcl
resource "azurerm_management_group" "example" {
  display_name = "Some Management Group"
}

resource "azurerm_policy_definition" "example" {
  name                = "only-deploy-in-westeurope"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "my-policy-definition"
  management_group_id = azurerm_management_group.example.id

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

resource "azurerm_management_group_policy_assignment" "example" {
  name                 = "example-policy"
  policy_definition_id = azurerm_policy_definition.example.id
  management_group_id  = azurerm_management_group.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `management_group_id` - (Required) The ID of the Management Group. Changing this forces a new Policy Assignment to be created.

* `name` - (Required) The name which should be used for this Policy Assignment. Changing this forces a new Policy Assignment to be created.

* `policy_definition_id` - (Required) The ID of the Policy Definition or Policy Definition Set. Changing this forces a new Policy Assignment to be created.

---

* `description` - (Optional) A description which should be used for this Policy Assignment.

* `display_name` - (Optional) The Display Name for this Policy Assignment.

* `enforce` - (Optional) Specifies if this Policy should be enforced or not?

* `identity` - (Optional) An `identity` block as defined below.

-> **Note:** The `location` field must also be specified when `identity` is specified.

* `location` - (Optional) The Azure Region where the Policy Assignment should exist. Changing this forces a new Policy Assignment to be created.

* `metadata` - (Optional) A JSON mapping of any Metadata for this Policy.

* `non_compliance_message` - (Optional) One or more `non_compliance_message` blocks as defined below.

* `not_scopes` - (Optional) Specifies a list of Resource Scopes (for example a Subscription, or a Resource Group) within this Management Group which are excluded from this Policy.

* `parameters` - (Optional) A JSON mapping of any Parameters for this Policy.

---

A `identity` block supports the following:

* `type` - (Optional) The Type of Managed Identity which should be added to this Policy Definition. Possible values are `SystemAssigned` and `UserAssigned`.

* `identity_ids` - (Optional) A list of User Managed Identity IDs which should be assigned to the Policy Definition.

~> **NOTE:** This is required when `type` is set to `UserAssigned`.

---

A `non_compliance_message` block supports the following:

* `content` - (Required) The non-compliance message text. When assigning policy sets (initiatives), unless `policy_definition_reference_id` is specified then this message will be the default for all policies.

* `policy_definition_reference_id` - (Optional) When assigning policy sets (initiatives), this is the ID of the policy definition that the non-compliance message applies to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Management Group Policy Assignment.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID of the Policy Assignment for this Management Group.

* `tenant_id` - The Tenant ID of the Policy Assignment for this Management Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Assignment for this Management Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Assignment for this Management Group.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Assignment for this Management Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Assignment for this Management Group.

## Import

Management Group Policy Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_management_group_policy_assignment.example /providers/Microsoft.Management/managementGroups/group1/providers/Microsoft.Authorization/policyAssignments/assignment1
```
