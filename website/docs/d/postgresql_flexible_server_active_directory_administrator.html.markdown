---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_postgresql_flexible_server_active_directory_administrator"
description: |-
  Gets information about an existing PostgreSQL Flexible Server Active Directory Administrator.
---

# Data Source: azurerm_postgresql_flexible_server_active_directory_administrator

Use this data source to access information about an existing PostgreSQL Flexible Server Active Directory Administrator.

## Example Usage

```hcl
data "azurerm_postgresql_flexible_server_active_directory_administrator" "example" {
  server_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/myserver1"
  object_id = "00000000-0000-0000-0000-000000000000"
}

output "principal_name" {
  value = data.azurerm_postgresql_flexible_server_active_directory_administrator.example.principal_name
}
```

## Arguments Reference

The following arguments are supported:

* `server_id` - (Required) The ID of the PostgreSQL Flexible Server on which to set the administrator.

* `object_id` - (Required) The object ID of a user, service principal or security group in the Azure Active Directory tenant set as the Flexible Server Admin.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL Flexible Server Active Directory Administrator.

* `principal_name` - The name of Azure Active Directory principal.

* `principal_type` - The type of Azure Active Directory principal.

* `tenant_id` - The Azure Tenant ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Flexible Server Active Directory Administrator.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DBforPostgreSQL` - 2024-08-01
