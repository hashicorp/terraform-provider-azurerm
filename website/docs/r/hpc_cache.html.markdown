---
subcategory: "Storage Cache"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hpc_cache"
description: |-
  Manages a HPC Cache.
---

# azurerm_hpc_cache

Manages a HPC Cache.

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
  cache_size          = 3
  subnet_id           = "${azurerm_subnet.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the HPC Cache. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which to create the HPC Cache. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure Region where the HPC Cache should be created. Changing this forces a new resource to be created.

~> **Please Note**: Not all locations support this resource. Some are `East US`, `East US 2`, `West US 2`, `North Europe`, and `West Europe`. 

* `cache_size` - (Required) The size of the HPC Cache, in GB. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet for the HPC Cache. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The `id` of the HPC Cache.

## Import

HPC Cache can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hpc_cache.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName/providers/Microsoft.StorageCache/caches/cacheName
```
