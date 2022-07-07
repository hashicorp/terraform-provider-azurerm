---
subcategory: "Disks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_pool_iscsi_target_lun"
description: |-
  Manages an iSCSI Target Lun(Logic Unit Number).
---

# azurerm_disk_pool_iscsi_target_lun

Manages an iSCSI Target lun.

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

resource "azurerm_disk_pool_iscsi_target" "example" {
  depends_on    = [azurerm_disk_pool_managed_disk_attachment.example]
  name          = "example"
  acl_mode      = "Dynamic"
  disks_pool_id = azurerm_disk_pool.example.id
  target_iqn    = "iqn.2021-11.com.microsoft:test"
}

resource "azurerm_disk_pool_iscsi_target_lun" "example" {
  iscsi_target_id                      = azurerm_disk_pool_iscsi_target.example.id
  disk_pool_managed_disk_attachment_id = azurerm_disk_pool_managed_disk_attachment.example.id
  name                                 = "example-disk"
}
```

## Arguments Reference

The following arguments are supported:

* `iscsi_target_id` - (Required) The ID of the iSCSI Target. Changing this forces a new iSCSI Target LUN to be created.

* `disk_pool_managed_disk_attachment_id` - (Required) The ID of the `azurerm_disk_pool_managed_disk_attachment`. Changing this forces a new iSCSI Target LUN to be created.

* `name` - (Required) User defined name for iSCSI LUN. Supported characters include uppercase letters, lowercase letters, numbers, periods, underscores or hyphens. Name should end with an alphanumeric character. The length must be between `1` and `90`. Changing this forces a new iSCSI Target LUN to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the iSCSI Target LUN.

* `lun` - The Logical Unit Number of the iSCSI Target LUN.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the iSCSI Target LUN.
* `read` - (Defaults to 5 minutes) Used when retrieving the iSCSI Target LUN.
* `delete` - (Defaults to 30 minutes) Used when deleting the iSCSI Target LUN.

## Import

iSCSI Target Luns can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_disk_pool_iscsi_target_lun.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.StoragePool/diskPools/diskPoolValue/iscsiTargets/iscsiTargetValue/lun|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/disks/disk1
```
