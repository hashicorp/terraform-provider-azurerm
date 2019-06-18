---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_container"
sidebar_current: "docs-azurerm-resource-storage-container"
description: |-
  Manages a Azure Storage Container.
---

# azurerm_storage_container

Manage an Azure Storage Container.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "westus"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststorageaccount"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "westus"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the storage container. Must be unique within the storage service the container is located.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the storage container. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) Specifies the storage account in which to create the storage container.
 Changing this forces a new resource to be created.

* `container_access_type` - (Optional) The 'interface' for access the container provides. Can be either `blob`, `container` or `private`. Defaults to `private`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Container.
* `properties` - Key-value definition of additional properties associated to the storage container

## Import

Storage Containers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_container.container1 https://example.blob.core.windows.net/container
```
