---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_policy_assignment"
description: |-
  Gets information about an existing Policy Assignment.
---

# Data Source: azurerm_policy_assignment

Use this data source to access information about an existing Policy Assignment.

## Example Usage

```hcl
data "azurerm_policy_assignment" "example" {
  name     = "existing"
  scope_id = data.azurerm_resource_group.example.id
}

output "id" {
  value = data.azurerm_policy_assignment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Policy Assignment. Changing this forces a new Policy Assignment to be created.

* `scope_id` - (Required) The ID of the scope this Policy Assignment is assigned to. The `scope_id` can be a subscription id, a resource group id, a management group id, or an ID of any resource that is assigned with a policy. Changing this forces a new Policy Assignment to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Policy Assignment.

* `description` - The description of this Policy Assignment.

* `display_name` - The display name of this Policy Assignment.

* `enforce` - Whether this Policy is enforced or not?

* `identity` - A `identity` block as defined below.

* `location` - The Azure Region where the Policy Assignment exists.

* `metadata` - A JSON mapping of any Metadata for this Policy.

* `non_compliance_message` - A `non_compliance_message` block as defined below.

* `not_scopes` - A `not_scopes` block as defined below.

* `parameters` - A JSON mapping of any Parameters for this Policy.

* `policy_definition_id` - The ID of the assigned Policy Definition.

---

A `identity` block exports the following:

* `identity_ids` - A `identity_ids` block as defined below.

* `principal_id` - The Principal ID of the Policy Assignment for this Resource.

* `tenant_id` - The Tenant ID of the Policy Assignment for this Resource.

* `type` - The Type of Managed Identity which is added to this Policy Assignment.

---

A `non_compliance_message` block exports the following:

* `content` - The non-compliance message text.

* `policy_definition_reference_id` - The ID of the Policy Definition that the non-compliance message applies to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Assignment.
