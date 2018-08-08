---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_container"
sidebar_current: "docs-azurerm-resource-storage-container"
description: |-
  Manages a Storage Container within a Storage Account.
---

# azurerm_storage_container

Manages a Storage Container within a Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}§

resource "azurerm_storage_account" "example" {
  # ...
}

resource "azurerm_storage_container" "example" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  storage_account_name  = "${azurerm_storage_account.example.name}"
  container_access_type = "private"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Storage Container. Changing this forces a new resource to be created.

-> **NOTE:** The `name` must be unique within the Storage Account where the Container is located.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the storage container. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) Specifies the storage account in which to create the storage container.
 Changing this forces a new resource to be created.

* `container_access_type` - (Optional) The 'interface' for access the container provides. Can be either `blob`, `container` or `private`. Defaults to `private`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The storage container Resource ID.
* `properties` - Key-value definition of additional properties associated to the storage container
