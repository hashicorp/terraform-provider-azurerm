---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_table"
sidebar_current: "docs-azurerm-resource-storage-table"
description: |-
  Manages a Storage Table within a Storage Account.
---

# azurerm_storage_table

Manages a Storage Table within a Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_storage_account" "example" {
  # ...
}

resource "azurerm_storage_table" "example" {
  name                 = "example-table"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  storage_account_name = "${azurerm_storage_account.example.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Storage Table. Changing this forces a new resource to be created.

-> **NOTE:** The `name` must be unique within the Storage Account where the Table is located.

* `resource_group_name` - (Required) The name of the resource group in which to create the storage table. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) Specifies the storage account in which to create the storage table. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The storage table Resource ID.
