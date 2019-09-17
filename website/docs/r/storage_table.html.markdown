---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_table"
sidebar_current: "docs-azurerm-resource-storage-table-x"
description: |-
  Manages a Table within an Azure Storage Account.
---

# azurerm_storage_table

Manages a Table within an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "azuretest"
  location = "West Europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "azureteststorage1"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "test" {
  name                 = "mysampletable"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the storage table. Must be unique within the storage account the table is located.

* `storage_account_name` - (Required) Specifies the storage account in which to create the storage table.
 Changing this forces a new resource to be created.

* `resource_group_name` - (Optional / **Deprecated**) The name of the resource group in which to create the storage table.

* `acl` - (Optional) One or more `acl` blocks as defined below.

---

A `acl` block supports the following:

* `id` - (Required) The ID which should be used for this Shared Identifier.

* `access_policy` - (Required) An `access_policy` block as defined below.

---

A `access_policy` block supports the following: 

* `expiry` - (Required) The ISO8061 UTC time at which this Access Policy should be valid until.

* `permissions` - (Required) The permissions which should associated with this Shared Identifier.

* `start` - (Required) The ISO8061 UTC time at which this Access Policy should be valid from.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Table within the Storage Account.

## Import

Table's within a Storage Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_table.table1 https://example.table.core.windows.net/table1
```
