---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_container"
sidebar_current: "docs-azurerm-resource-storage-container"
description: |-
  Manages a Container within an Azure Storage Account.
---

# azurerm_storage_container

Manages a Container within an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "examplestoraccount"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
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

* `name` - (Required) The name of the Container which should be created within the Storage Account.

* `storage_account_name` - (Required) The name of the Storage Account where the Container should be created.

* `container_access_type` - (Optional) The Access Level configured for this Container. Possible values are `blob`, `container` or `private`. Defaults to `private`.

* `metadata` - (Optional) A mapping of MetaData for this Container.

* `resource_group_name` - (Optional / **Deprecated**) The name of the resource group in which to create the storage container. This field is no longer used and will be removed in 2.0. 

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Container.

* `has_immutability_policy` - Is there an Immutability Policy configured on this Storage Container?

* `has_legal_hold` - Is there a Legal Hold configured on this Storage Container?

* `properties` - (**Deprecated**) Key-value definition of additional properties associated to the Storage Container

## Import

Storage Containers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_container.container1 https://example.blob.core.windows.net/container
```
