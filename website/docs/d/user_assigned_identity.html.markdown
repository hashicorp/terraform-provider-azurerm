---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azure_user_assigned_identity"
description: |-
  Gets information about an existing User Assigned Identity.

---

# Data Source: azurerm_user_assigned_identity

Use this data source to access information about an existing User Assigned Identity.

## Example Usage (reference an existing)

```hcl
data "azurerm_user_assigned_identity" "example" {
  name                = "name_of_user_assigned_identity"
  resource_group_name = "name_of_resource_group"
}

output "uai_client_id" {
  value = data.azurerm_user_assigned_identity.example.client_id
}

output "uai_principal_id" {
  value = data.azurerm_user_assigned_identity.example.principal_id
}

output "uai_tenant_id" {
  value = data.azurerm_user_assigned_identity.example.tenant_id
}
```

## Argument Reference

* `name` -  The name of the User Assigned Identity.
* `resource_group_name` - The name of the Resource Group in which the User Assigned Identity exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the User Assigned Identity.
* `location` - The Azure location where the User Assigned Identity exists.
* `principal_id` - The Service Principal ID of the User Assigned Identity.
* `client_id` - The Client ID of the User Assigned Identity.
* `tenant_id` - The Tenant ID of the User Assigned Identity.
* `tags` - A mapping of tags assigned to the User Assigned Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the User Assigned Identity.
