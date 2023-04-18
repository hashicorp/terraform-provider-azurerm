---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_target_endpoint"
description: |-
  Manages a Storage Mover Target Endpoint.
---

# azurerm_storage_mover_target_endpoint

Manages a Storage Mover Target Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                            = "examplestr"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "example" {
  name                  = "example-sc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "blob"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
}

resource "azurerm_storage_mover_target_endpoint" "example" {
  name                   = "example-se"
  storage_mover_id       = azurerm_storage_mover.example.id
  storage_account_id     = azurerm_storage_account.example.id
  storage_container_name = azurerm_storage_container.example.name
  description            = "Example Storage Container Endpoint Description"
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover Target Endpoint. Changing this forces a new resource to be created.

* `storage_mover_id` - (Required) Specifies the ID of the storage mover for this Storage Mover Target Endpoint. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) Specifies the ID of the storage account for this Storage Mover Target Endpoint. Changing this forces a new resource to be created.

* `storage_container_name` - (Required) Specifies the name of the storage blob container for this Storage Mover Target Endpoint. Changing this forces a new resource to be created.

* `description` - (Optional) Specifies a description for the Storage Mover Target Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover Target Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover Target Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover Target Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover Target Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover Target Endpoint.

## Import

Storage Mover Target Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_target_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/endpoints/endpoint1
```
