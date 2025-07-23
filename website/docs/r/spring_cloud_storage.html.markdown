---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_storage"
description: |-
  Manages a Spring Cloud Storage.
---

# azurerm_spring_cloud_storage

Manages a Spring Cloud Storage.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_storage` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_spring_cloud_storage" "example" {
  name                    = "example"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
  storage_account_name    = azurerm_storage_account.example.name
  storage_account_key     = azurerm_storage_account.example.primary_access_key
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Storage. Changing this forces a new Spring Cloud Storage to be created.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service where the Spring Cloud Storage should exist. Changing this forces a new Spring Cloud Storage to be created.

* `storage_account_key` - (Required) The access key of the Azure Storage Account.

* `storage_account_name` - (Required) The account name of the Azure Storage Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Storage.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Storage.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Storage.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Storage.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Storage.

## Import

Spring Cloud Storages can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_storage.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/storages/storage1
```
