---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server_extended_auditing_policy"
description: |-
  Manages a Ms Sql Server Extended Auditing Policy.
---

# azurerm_mssql_server_extended_auditing_policy

Manages a Ms Sql Server Extended Auditing Policy.

~> **NOTE:** The Server Extended Auditing Policy Can be set inline here as well as with the [mssql_server_extended_auditing_policy resource](mssql_server_extended_auditing_policy.html) resource. You can only use one or the other and using both will cause a conflict.

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

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_server_extended_auditing_policy" "example" {
  server_id                               = azurerm_mssql_server.example.id
  storage_endpoint                        = azurerm_storage_account.example.primary_blob_endpoint
  storage_account_access_key              = azurerm_storage_account.example.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days                       = 6
}
```

## Arguments Reference

The following arguments are supported:

* `server_id` - (Required) The ID of the sql server to set the extended auditing policy. Changing this forces a new resource to be created.

* `storage_endpoint` - (Required) The blob storage endpoint (e.g. https://MyAccount.blob.core.windows.net). This blob storage will hold all extended auditing logs.

---

* `retention_in_days` - (Optional) The number of days to retain logs for in the storage account.

* `storage_account_access_key` - (Optional) The access key to use for the auditing storage account.

* `storage_account_access_key_is_secondary` - (Optional) Is `storage_account_access_key` value the storage's secondary key?

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Ms Sql Server Extended Auditing Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Ms Sql Server Extended Auditing Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Ms Sql Server Extended Auditing Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Ms Sql Server Extended Auditing Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Ms Sql Server Extended Auditing Policy.

## Import

Ms Sql Server Extended Auditing Policys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_server_extended_auditing_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/servers/sqlServer1/extendedAuditingSettings/default
```
