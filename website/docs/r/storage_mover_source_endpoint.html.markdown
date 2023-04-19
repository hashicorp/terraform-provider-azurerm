---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_source_endpoint"
description: |-
  Manages a Storage Mover Source Endpoint.
---

# azurerm_storage_mover_source_endpoint

Manages a Storage Mover Source Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
}

resource "azurerm_storage_mover_source_endpoint" "example" {
  name             = "example-se"
  storage_mover_id = azurerm_storage_mover.example.id
  export           = "/"
  host             = "192.168.0.1"
  nfs_version      = "NFSv3"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover Source Endpoint. Changing this forces a new resource to be created.

* `storage_mover_id` - (Required) Specifies the ID of the Storage Mover for this Storage Mover Source Endpoint. Changing this forces a new resource to be created.

* `host` - (Required) Specifies the host name or IP address of the server exporting the file system. Changing this forces a new resource to be created.

* `export` - (Optional) Specifies the directory being exported from the server. Changing this forces a new resource to be created.

* `nfs_version` - (Optional) Specifies the NFS protocol version. Possible values are `NFSauto`, `NFSv3` and `NFSv4`. Defaults to `NFSauto`. Changing this forces a new resource to be created.

* `description` - (Optional) Specifies a description for the Storage Mover Source Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover Source Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover Source Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover Source Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover Source Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover Source Endpoint.

## Import

Storage Mover Source Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_source_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/endpoints/endpoint1
```
