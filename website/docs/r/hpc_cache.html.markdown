---
subcategory: "Storage Cache"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hpc_cache"
description: |-
  Manages a HPC Cache.
---

# azurerm_hpc_cache

Manages a HPC Cache.

~> **Note**: During the first several months of the GA release, a request must be made to the Azure HPC Cache team to add your subscription to the access list before it can be used to create a cache instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "examplerg"
  location = "eastus"
}

resource "azurerm_virtual_network" "example" {
  name                = "examplevn"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_subnet" "example" {
  name                 = "examplesubnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_hpc_cache" "example" {
  name                = "examplehpccache"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  cache_size          = 3072
  subnet_id           = "${azurerm_subnet.example.id}"
  sku_name            = "Standard_2G"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the HPC Cache. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which to create the HPC Cache. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure Region where the HPC Cache should be created. Changing this forces a new resource to be created.

~> **Note**: Not all locations support this resource. Supported regions are `East US`, `East US 2`, `West US 2`, `North Europe`, `West Europe`, `Southeast Asia`, `Korea Central` and `Sydney`. 

* `cache_size` - (Required) The size of the HPC Cache, in GB. Possible values are `3072`, `6144`, `12288`, `24576`, and `49152`. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet for the HPC Cache. Changing this forces a new resource to be created.

* `sku_name` - (Required) The SKU of HPC Chace to use. Possible values are `Standard_2G`, `Standard_4G` and `Standard_8G`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The `id` of the HPC Cache.

## Import

HPC Cache can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hpc_cache.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName/providers/Microsoft.StorageCache/caches/cacheName
```
