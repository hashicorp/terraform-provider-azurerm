---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_policy_assignment"
description: |-
  Manages a Policy Assignment to a Resource.
---

# azurerm_resource_policy_assignment

Manages a Policy Assignment to a Resource.

## Example Usage

```hcl
data "azurerm_virtual_network" "example" {
  name                = "production"
  resource_group_name = "networking"
}

resource "azurerm_policy_definition" "example" {
  name        = "only-deploy-in-westeurope"
  policy_type = "Custom"
  mode        = "All"

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

resource "azurerm_resource_policy_assignment" "example" {
  name                 = "example-policy-assignment"
  resource_id          = data.azurerm_virtual_network.example.id
  policy_definition_id = azurerm_policy_definition.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Policy Assignment. Changing this forces a new Resource Policy Assignment to be created.

* `policy_definition_id` - (Required) The ID of the Policy Definition or Policy Definition Set. Changing this forces a new Policy Assignment to be created.

* `resource_id` - (Required) The ID of the Resource (or Resource Scope) where this should be applied. Changing this forces a new Resource Policy Assignment to be created.

~> To create a Policy Assignment at a Management Group use the `azurerm_management_group_policy_assignment` resource, for a Resource Group use the `azurerm_resource_group_policy_assignment` and for a Subscription use the `azurerm_subscription_policy_assignment` resource.

---

* `description` - (Optional) A description which should be used for this Policy Assignment.

* `display_name` - (Optional) The Display Name for this Policy Assignment.

* `enforce` - (Optional) Specifies if this Policy should be enforced or not?

* `identity` - (Optional) A `identity` block as defined below.

-> **Note:** The `location` field must also be specified when `identity` is specified.

* `location` - (Optional) The Azure Region where the Policy Assignment should exist. Changing this forces a new Policy Assignment to be created.

* `metadata` - (Optional) A JSON mapping of any Metadata for this Policy.

* `not_scopes` - (Optional) Specifies a list of Resource Scopes (for example a Subscription, or a Resource Group) within this Management Group which are excluded from this Policy.

* `parameters` - (Optional) A JSON mapping of any Parameters for this Policy. Changing this forces a new Management Group Policy Assignment to be created.

---

A `identity` block supports the following:

* `type` - (Optional) The Type of Managed Identity which should be added to this Policy Definition. The only possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Policy Assignment.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID of the Policy Assignment for this Resource.

* `tenant_id` - The Tenant ID of the Policy Assignment for this Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Assignment for this Resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Assignment for this Resource.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Assignment for this Resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Assignment for this Resource.

## Import

Resource Policy Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_policy_assignment.example "{resource}/providers/Microsoft.Authorization/policyAssignments/assignment1"
```

where `{resource}` is a Resource ID in the form `/subscriptions/00000000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualNetworks/network1`.
