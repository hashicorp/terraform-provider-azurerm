---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server_dns_alias"
description: |-
  Manages a MS SQL Server DNS Alias.
---

# azurerm_mssql_server_dns_alias

Manages a MS SQL Server DNS Alias.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mssql_server" "example" {
  name                         = "example-sqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "AdminPassword123!"
}

resource "azurerm_mssql_server_dns_alias" "example" {
  name            = "example-dns-alias"
  mssql_server_id = azurerm_mssql_server.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `mssql_server_id` - (Required) The ID of the mssql server. Changing this forces a new MSSQL Server DNS Alias to be created.

* `name` - (Required) The name which should be used for this MSSQL Server DNS Alias. Changing this forces a new MSSQL Server DNS Alias to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MSSQL Server DNS Alias.

* `dns_record` - The fully qualified DNS record for alias.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MSSQL Server DNS Alias.
* `read` - (Defaults to 5 minutes) Used when retrieving the MSSQL Server DNS Alias.
* `delete` - (Defaults to 10 minutes) Used when deleting the MSSQL Server DNS Alias.

## Import

MSSQL Server DNS Aliass can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_server_dns_alias.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/dnsAliases/default
```
