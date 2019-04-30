---
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
data "azurerm_user_assigned_identity" "test" {
  name                = "name_of_user_assigned_identity"
  resource_group_name = "name_of_resource_group"
}

output "uai_client_id" {
  value = "${data.azurerm_user_assigned_identity.test.client_id}"
}

output "uai_principal_id" {
  value = "${data.azurerm_user_assigned_identity.test.principal_id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the user assigned identity.
* `resource_group_name` - (Required) Specifies the name of the resource group.

## Attributes Reference

The following attributes are exported:

* `id` - The user assigned identity ID.
* `principal_id` - Service Principal ID associated with the user assigned identity.
* `client_id` - Client ID associated with the user assigned identity.
