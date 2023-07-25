---
subcategory: "Disks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_pool_managed_disk_attachment"
description: |-
  Manages a Disk Pool Managed Disk Attachment.
---

# azurerm_disk_pool_managed_disk_attachment

Manages a Disk Pool Managed Disk Attachment.

!> **Note:** Azure are officially [halting](https://learn.microsoft.com/en-us/azure/azure-vmware/attach-disk-pools-to-azure-vmware-solution-hosts?tabs=azure-cli) the preview of Azure Disk Pools, and it **will not** be made generally available. New customers will not be able to register the Microsoft.StoragePool resource provider on their subscription and deploy new Disk Pools. Existing subscriptions registered with Microsoft.StoragePool may continue to deploy and manage disk pools for the time being.

~> **Note:** Must be either a premium SSD, standard SSD, or an ultra disk in the same region and availability zone as the disk pool.

~> **Note:** Ultra disks must have a disk sector size of 512 bytes.

~> **Note:** Must be a shared disk, with a maxShares value of two or greater.

~> **Note:** You must provide the StoragePool resource provider RBAC permissions to the disks that will be added to the disk pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/24"]
  delegation {
    name = "diskspool"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/read"]
      name    = "Microsoft.StoragePool/diskPools"
    }
  }
}

resource "azurerm_disk_pool" "example" {
  name                = "example-pool"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  subnet_id           = azurerm_subnet.example.id
  zones               = ["1"]
  sku_name            = "Basic_B1"
}

resource "azurerm_managed_disk" "example" {
  name                 = "example-disk"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  create_option        = "Empty"
  storage_account_type = "Premium_LRS"
  disk_size_gb         = 4
  max_shares           = 2
  zone                 = "1"
}

data "azuread_service_principal" "example" {
  display_name = "StoragePool Resource Provider"
}

locals {
  roles = ["Disk Pool Operator", "Virtual Machine Contributor"]
}

resource "azurerm_role_assignment" "example" {
  count                = length(local.roles)
  principal_id         = data.azuread_service_principal.example.id
  role_definition_name = local.roles[count.index]
  scope                = azurerm_managed_disk.example.id
}

resource "azurerm_disk_pool_managed_disk_attachment" "example" {
  depends_on      = [azurerm_role_assignment.example]
  disk_pool_id    = azurerm_disk_pool.example.id
  managed_disk_id = azurerm_managed_disk.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `disk_pool_id` - (Required) The ID of the Disk Pool. Changing this forces a new Disk Pool Managed Disk Attachment to be created.

* `managed_disk_id` - (Required) The ID of the Managed Disk. Changing this forces a new Disks Pool Managed Disk Attachment to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Disk Pool Managed Disk Attachment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Disks Pool Managed Disk Attachment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Disks Pool Managed Disk Attachment.
* `delete` - (Defaults to 60 minutes) Used when deleting the Disks Pool Managed Disk Attachment.

## Import

Disks Pool Managed Disk Attachments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_disk_pool_managed_disk_attachment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StoragePool/diskPools/storagePool1/managedDisks|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/disks/disk1
```
