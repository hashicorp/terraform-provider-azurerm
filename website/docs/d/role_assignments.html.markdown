---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_role_assignments"
description: |-
  Gets information about existing Role Assignments.
---

# Data Source: azurerm_role_assignments

Use this data source to access information about existing Role Assignments.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

data "azurerm_role_assignments" "example" {
  scope = azurerm_resource_group.example.id
}

output "id" {
  value = data.azurerm_role_assignments.example.role_assignments
}
```

## Arguments Reference

The following arguments are supported:

* `scope` - (Required) The scope at which to list Role Assignments.

---

* `limit_at_scope` - (Optional) Whether to limit the result exactly at the specified scope and not above or below it. Defaults to `false`.

* `principal_id` - (Optional) The principal ID to filter the list of Role Assignments against.

* `tenant_id` - (Optional) The tenant ID for cross-tenant requests.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of this Role Assignments data source.

* `role_assignments` - A `role_assignments` block as defined below.

---

A `role_assignments` block exports the following:

* `condition` - The condition that limits the resource the role can be assigned to.

* `condition_version` - The version of the condition.

* `delegated_managed_identity_resource_id` - The ID of the delegated managed identity resource.

* `description` - The description for this Role Assignment.

* `principal_id` - The principal ID.

* `principal_type` - The type of the `principal_id`.

* `role_assignment_id` - The ID of the Role Assignment.

* `role_assignment_name` - The name of the Role Assignment.

* `role_assignment_scope` - The scope of the Role Assignment.

* `role_definition_id` - The ID of the Role Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Role Assignments.
