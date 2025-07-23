---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_dev_center_project_environment_type"
description: |-
  Gets information about an existing Dev Center Project Environment Type.
---

# Data Source: azurerm_dev_center_project_environment_type

Use this data source to access information about an existing Dev Center Project Environment Type.

## Example Usage

```hcl
data "azurerm_dev_center_project_environment_type" "example" {
  name                  = azurerm_dev_center_project_environment_type.example.name
  dev_center_project_id = azurerm_dev_center_project_environment_type.example.dev_center_project_id
}

output "id" {
  value = data.azurerm_dev_center_project_environment_type.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Dev Center Project Environment Type.

* `dev_center_project_id` - (Required) The ID of the associated Dev Center Project.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Project Environment Type.

* `location` - The Azure Region where the Dev Center Project Environment Type exists.

* `deployment_target_id` - The ID of the subscription that the Environment Type is mapped to.

* `identity` - An `identity` block as defined below.

* `creator_role_assignment_roles` - A list of roles assigned to the environment creator.

* `user_role_assignment` - A `user_role_assignment` block as defined below.

* `tags` - A mapping of tags assigned to the Dev Center Project Environment Type.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Dev Center Project Environment Type.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Dev Center Project Environment Type.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Dev Center Project Environment Type.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Dev Center Project Environment Type.

---

A `user_role_assignment` block supports the following:

* `user_id` - The user object ID that is assigned roles.

* `roles` - A list of roles to assign to the `user_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Project Environment Type.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DevCenter`: 2025-02-01
