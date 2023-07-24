---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hpc_cache_aml_file_system"
description: |-
  Manages a HPC Cache AML File System.
---

# azurerm_hpc_cache_aml_file_system

Manages a HPC Cache AML File System.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_hpc_cache_aml_file_system" "example" {
  name                   = "example-amlfs"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  sku_name               = "AMLFS-Durable-Premium-250"
  subnet_id              = azurerm_subnet.example.id
  storage_capacity_in_tb = 8
  zones                  = ["2"]

  maintenance_window {
    day_of_week     = "Friday"
    time_of_day_utc = "22:00"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this HPC Cache AML File System. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the HPC Cache AML File System should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the HPC Cache AML File System should exist. Changing this forces a new resource to be created.

* `maintenance_window` - (Required) A `maintenance_window` block as defined below.

* `sku_name` - (Required) The SKU name for the HPC Cache AML File System. Changing this forces a new resource to be created.

* `storage_capacity_in_tb` - (Required) The size of the HPC Cache AML File System in TiB. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The subnet used for managing the AML file system and for client-facing operations. This subnet should have at least a /24 subnet mask within the VNET's address space. Changing this forces a new resource to be created.

* `zones` - (Required) The availability zones for the HPC Cache AML File System. Changing this forces a new resource to be created.

* `hsm_setting` - (Optional) A `hsm_setting` block as defined below. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below. Changing this forces a new resource to be created.

* `key_encryption_key` - (Optional) A `key_encryption_key` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the HPC Cache AML File System.

---

A `maintenance_window` block supports the following:

* `day_of_week` - (Required) The day of the week on which the maintenance window will occur.

* `time_of_day_utc` - (Required) The time of day (in UTC) to start the maintenance window.

---

A `hsm_setting` block supports the following:

* `container_id` - (Required) The resource ID of storage container used for hydrating the namespace and archiving from the namespace. Changing this forces a new resource to be created.

* `logging_container_id` - (Required) The resource ID of storage container used for logging events and errors. Changing this forces a new resource to be created.

* `import_prefix` - (Optional) Only blobs in the non-logging container that start with this path/prefix get hydrated into the cluster namespace. Changing this forces a new resource to be created.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that should be configured on this HPC Cache AML File System. Only possible value is `UserAssigned`. Changing this forces a new resource to be created.

* `identity_ids` - (Required) A list of User Assigned Managed Identity IDs to be assigned to this HPC Cache AML File System. Changing this forces a new resource to be created.

---

A `key_encryption_key` block supports the following:

* `key_url` - (Required) The URL to the Key Vault Key used as the Key Encryption Key. This can be found as `id` on the `azurerm_key_vault_key` resource.

* `source_vault_id` - (Required) The ID of the source Key Vault. This can be found as `id` on the `azurerm_key_vault` resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the HPC Cache AML File System.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the HPC Cache AML File System.
* `read` - (Defaults to 5 minutes) Used when retrieving the HPC Cache AML File System.
* `update` - (Defaults to 30 minutes) Used when updating the HPC Cache AML File System.
* `delete` - (Defaults to 30 minutes) Used when deleting the HPC Cache AML File System.

## Import

HPC Cache AML File Systems can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hpc_cache_aml_file_system.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageCache/amlFilesystems/amlFilesystem1
```
