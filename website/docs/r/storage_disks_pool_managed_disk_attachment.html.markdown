---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_disks_pool_managed_disk_attachment"
description: |-
  Manages a Disks Pool Managed Disk Attachment.
---

# azurerm_storage_disks_pool_managed_disk_attachment

Manages a Disks Pool Managed Disk Attachment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "Central US"
}

resource "azurerm_virtual_network" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example"
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

resource "azurerm_storage_disks_pool" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  subnet_id           = azurerm_subnet.example.id
  availability_zones  = ["1"]
  sku_name            = "Basic_B1"
}

resource "azurerm_managed_disk" "example" {
  name                 = "example"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  create_option        = "Empty"
  storage_account_type = "Premium_LRS"
  disk_size_gb         = 4
  max_shares           = 2
  zones                = ["1"]
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

resource "azurerm_storage_disks_pool_managed_disk_attachment" "example" {
  depends_on      = [azurerm_role_assignment.example]
  disks_pool_id   = azurerm_storage_disks_pool.example.id
  managed_disk_id = azurerm_managed_disk.test.id
}
```

## Arguments Reference

The following arguments are supported:

* `disks_pool_id` - (Required) The ID of the Disks Pool. Changing this forces a new Disks Pool Managed Disk Attachment to be created.

* `managed_disk_id` - (Required) The ID of the Managed Disk. Changing this forces a new Disks Pool Managed Disk Attachment to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Disks Pool Managed Disk Attachment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Disks Pool Managed Disk Attachment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Disks Pool Managed Disk Attachment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Disks Pool Managed Disk Attachment.

## Import

Disks Pool Managed Disk Attachments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_disks_pool_managed_disk_attachment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StoragePool/diskPools/storagePool1/managedDisks|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/disks/disk1
```
