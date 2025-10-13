---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_table"
description: |-
  Manages a Table within an Azure Storage Account.
---

# azurerm_storage_table

Manages a Table within an Azure Storage Account.

~> **Note:** Shared Key authentication will always be used for this resource, as AzureAD authentication is not supported when setting or retrieving ACLs for Tables.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "azuretest"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "azureteststorage1"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "example" {
  name                 = "mysampletable"
  storage_account_name = azurerm_storage_account.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the storage table. Only Alphanumeric characters allowed, starting with a letter. Must be unique within the storage account the table is located. Changing this forces a new resource to be created.

* `storage_account_name` - (Optional) The name of the Storage Account where the Storage Table should be created. Changing this forces a new resource to be created. This property is deprecated in favour of `storage_account_id`.

~> **Note:** Migrating from the deprecated `storage_account_name` to `storage_account_id` is supported without recreation. Any other change to either property will result in the resource being recreated.

* `storage_account_id` - (Optional) The name of the Storage Account where the Storage Table should be created. Changing this forces a new resource to be created.

~> **Note:** One of `storage_account_name` or `storage_account_id` must be specified. When specifying `storage_account_id` the resource will use the Resource Manager API, rather than the Data Plane API.

* `acl` - (Optional) One or more `acl` blocks as defined below.

---

A `acl` block supports the following:

* `id` - (Required) The ID which should be used for this Shared Identifier.

* `access_policy` - (Optional) An `access_policy` block as defined below.

---

A `access_policy` block supports the following:

* `expiry` - (Required) The ISO8061 UTC time at which this Access Policy should be valid until.

* `permissions` - (Required) The permissions which should associated with this Shared Identifier.

* `start` - (Required) The ISO8061 UTC time at which this Access Policy should be valid from.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Table within the Storage Account.

* `resource_manager_id` - The Resource Manager ID of this Storage Table.

* `url` - The data plane URL of the Storage Table in the format of `<storage table endpoint>/Tables('<table name>')`. E.g. `https://example.table.core.windows.net/Tables('mytable')"`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Table.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Table.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Table.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Table.

## Import

Table's within a Storage Account can be imported using the `resource id`, e.g.

If `storage_account_name` is used:

```shell
terraform import azurerm_storage_table.table1 "https://example.table.core.windows.net/Tables('replace-with-table-name')"
```

If `storage_account_id` is used:

```shell
terraform import azurerm_storage_table.table1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount/tableServices/default/tables

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Storage` - 2023-05-01
