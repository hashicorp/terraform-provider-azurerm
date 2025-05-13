---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hpc_cache_blob_target"
description: |-
  Manages a Blob Target within a HPC Cache.
---

# azurerm_hpc_cache_blob_target

Manages a Blob Target within a HPC Cache.

~> **Note:** By request of the service team the provider no longer automatically registering the `Microsoft.StorageCache` Resource Provider for this resource. To register it you can run `az provider register --namespace 'Microsoft.StorageCache'`.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "examplevn"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "examplesubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_hpc_cache" "example" {
  name                = "examplehpccache"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.example.id
  sku_name            = "Standard_2G"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorgaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                 = "examplestoragecontainer"
  storage_account_name = azurerm_storage_account.example.name
}

data "azuread_service_principal" "example" {
  display_name = "HPC Cache Resource Provider"
}

resource "azurerm_role_assignment" "example_storage_account_contrib" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Account Contributor"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_role_assignment" "example_storage_blob_data_contrib" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_hpc_cache_blob_target" "example" {
  name                 = "examplehpccblobtarget"
  resource_group_name  = azurerm_resource_group.example.name
  cache_name           = azurerm_hpc_cache.example.name
  storage_container_id = azurerm_storage_container.example.resource_manager_id
  namespace_path       = "/blob_storage"
}
```

## Argument Reference

The following arguments are supported:

* `cache_name` - (Required) The name HPC Cache, which the HPC Cache Blob Target will be added to. Changing this forces a new resource to be created.

* `name` - (Required) The name of the HPC Cache Blob Target. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which to create the HPC Cache Blob Target. Changing this forces a new resource to be created.

* `namespace_path` - (Required) The client-facing file path of the HPC Cache Blob Target.

* `storage_container_id` - (Required) The Resource Manager ID of the Storage Container used as the HPC Cache Blob Target. Changing this forces a new resource to be created.

-> **Note:** This is the Resource Manager ID of the Storage Container, rather than the regular ID - and can be accessed on the `azurerm_storage_container` Data Source/Resource as `resource_manager_id`.

* `access_policy_name` - (Optional) The name of the access policy applied to this target. Defaults to `default`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the HPC Cache Blob Target.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the HPC Cache Blob Target.
* `read` - (Defaults to 5 minutes) Used when retrieving the HPC Cache Blob Target.
* `update` - (Defaults to 30 minutes) Used when updating the HPC Cache Blob Target.
* `delete` - (Defaults to 30 minutes) Used when deleting the HPC Cache Blob Target.

## Import

Blob Targets within an HPC Cache can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hpc_cache_blob_target.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StorageCache/caches/cache1/storageTargets/target1
```
