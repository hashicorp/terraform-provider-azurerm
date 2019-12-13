---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azure_user_assigned_identity"
sidebar_current: "docs-azurerm-datasource-user-assigned-identity"
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
  value = "${data.azurerm_user_assigned_identity.example.client_id}"
}

output "uai_principal_id" {
  value = "${data.azurerm_user_assigned_identity.example.principal_id}"
}
```

## Argument Reference

* `name` - (Required)  The name of the User Assigned Identity.
* `resource_group_name` - (Required) The name of the Resource Group in which the User Assigned Identity exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the User Assigned Identity.
* `location` - The Azure location where the User Assigned Identity exists.
* `principal_id` - The Service Principal ID of the User Assigned Identity.
* `client_id` - The Client ID of the User Assigned Identity.
* `tags` - A mapping of tags assigned to the User Assigned Identity.
