---
subcategory: "Azure Managed Lustre File System"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_lustre_file_system"
description: |-
  Manages an Azure Managed Lustre File System.
---

# azurerm_managed_lustre_file_system

Manages an Azure Managed Lustre File System.

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

resource "azurerm_managed_lustre_file_system" "example" {
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

* `name` - (Required) The name which should be used for this Azure Managed Lustre File System. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Managed Lustre File System should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Managed Lustre File System should exist. Changing this forces a new resource to be created.

* `maintenance_window` - (Required) A `maintenance_window` block as defined below.

* `sku_name` - (Required) The SKU name for the Azure Managed Lustre File System. Possible values are `AMLFS-Durable-Premium-40`, `AMLFS-Durable-Premium-125`, `AMLFS-Durable-Premium-250` and `AMLFS-Durable-Premium-500`. Changing this forces a new resource to be created.

* `storage_capacity_in_tb` - (Required) The size of the Azure Managed Lustre File System in TiB. The valid values for this field are dependant on which `sku_name` has been defined in the configuration file. For more information on the valid values for this field please see the [product documentation](https://learn.microsoft.com/azure/azure-managed-lustre/create-file-system-resource-manager#file-system-type-and-size-options). Changing this forces a new resource to be created.


* `subnet_id` - (Required) The resource ID of the Subnet that is used for managing the Azure Managed Lustre file system and for client-facing operations. This subnet should have at least a /24 subnet mask within the Virtual Network's address space. Changing this forces a new resource to be created.

* `zones` - (Required) A list of availability zones for the Azure Managed Lustre File System. Changing this forces a new resource to be created.

* `hsm_setting` - (Optional) A `hsm_setting` block as defined below. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below. Changing this forces a new resource to be created.

* `encryption_key` - (Optional) An `encryption_key` block as defined below.

-> **Note:** Removing `encryption_key` forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Managed Lustre File System.

---

A `maintenance_window` block supports the following:

* `day_of_week` - (Required) The day of the week on which the maintenance window will occur. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` and `Saturday`.

* `time_of_day_in_utc` - (Required) The time of day (in UTC) to start the maintenance window.

---

A `hsm_setting` block supports the following:

* `container_id` - (Required) The resource ID of the storage container that is used for hydrating the namespace and archiving from the namespace. Changing this forces a new resource to be created.

* `logging_container_id` - (Required) The resource ID of the storage container that is used for logging events and errors. Changing this forces a new resource to be created.

* `import_prefix` - (Optional) The import prefix for the Azure Managed Lustre File System. Only blobs in the non-logging container that start with this path/prefix get hydrated into the cluster namespace. Changing this forces a new resource to be created.

-> **Note:** The roles `Contributor` and `Storage Blob Data Contributor` must be added to the Service Principal `HPC Cache Resource Provider` for the Storage Account. See [official docs]( https://learn.microsoft.com/en-us/azure/azure-managed-lustre/amlfs-prerequisites#access-roles-for-blob-integration) for more information.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that should be configured on this Azure Managed Lustre File System. Only possible value is `UserAssigned`. Changing this forces a new resource to be created.

* `identity_ids` - (Required) A list of User Assigned Managed Identity IDs to be assigned to this Azure Managed Lustre File System. Changing this forces a new resource to be created.

---

An `encryption_key` block supports the following:

* `key_url` - (Required) The URL to the Key Vault Key used as the Encryption Key. This can be found as `id` on the `azurerm_key_vault_key` resource.

* `source_vault_id` - (Required) The ID of the source Key Vault. This can be found as `id` on the `azurerm_key_vault` resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Managed Lustre File System.

* `mgs_address` - IP Address of Managed Lustre File System Services.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Managed Lustre File System.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Managed Lustre File System.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Managed Lustre File System.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Managed Lustre File System.

## Import

Azure Managed Lustre File Systems can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_lustre_file_system.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageCache/amlFilesystems/amlFilesystem1
```
