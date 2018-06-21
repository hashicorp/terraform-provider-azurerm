---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_table"
sidebar_current: "docs-azurerm-resource-storage-table"
description: |-
  Create a Azure Storage Table.
---

# azurerm_storage_table

Create an Azure Storage Table.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "azuretest"
  location = "westus"
}

resource "azurerm_storage_account" "test" {
  name                     = "azureteststorage1"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "westus"
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

* `resource_group_name` - (Required) The name of the resource group in which to
    create the storage table. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) Specifies the storage account in which to create the storage table.
 Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The storage table Resource ID.
