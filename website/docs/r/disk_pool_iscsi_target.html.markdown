---
subcategory: "Disks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_pool_iscsi_target"
description: |-
  Manages an iSCSI Target.
---

# azurerm_disk_pool_iscsi_target

Manages an iSCSI Target.

!> **Note:** Azure are officially [halting](https://learn.microsoft.com/en-us/azure/azure-vmware/attach-disk-pools-to-azure-vmware-solution-hosts?tabs=azure-cli) the preview of Azure Disk Pools, and it **will not** be made generally available. New customers will not be able to register the Microsoft.StoragePool resource provider on their subscription and deploy new Disk Pools. Existing subscriptions registered with Microsoft.StoragePool may continue to deploy and manage disk pools for the time being.

!> **Note:** Each Disk Pool can have a maximum of 1 iSCSI Target.

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
```

## Arguments Reference

The following arguments are supported:

* `acl_mode` - (Required) Mode for Target connectivity. The only supported value is `Dynamic` for now. Changing this forces a new iSCSI Target to be created.

* `disks_pool_id` - (Required) The ID of the Disk Pool. Changing this forces a new iSCSI Target to be created.

* `name` - (Required) The name of the iSCSI Target. The name can only contain lowercase letters, numbers, periods, or hyphens, and length should between [5-223]. Changing this forces a new iSCSI Target to be created.

* `target_iqn` - (Optional) ISCSI Target IQN (iSCSI Qualified Name); example: `iqn.2005-03.org.iscsi:server`. IQN should follow the format `iqn.yyyy-mm.<abc>.<pqr>[:xyz]`; supported characters include alphanumeric characters in lower case, hyphen, dot and colon, and the length should between `4` and `223`. Changing this forces a new iSCSI Target to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the iSCSI Target.

* `endpoints` - List of private IPv4 addresses to connect to the iSCSI Target.

* `port` - The port used by iSCSI Target portal group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the iSCSI Target.
* `read` - (Defaults to 5 minutes) Used when retrieving the iSCSI Target.
* `update` - (Defaults to 30 minutes) Used when updating the iSCSI Target.
* `delete` - (Defaults to 60 minutes) Used when deleting the iSCSI Target.

## Import

iSCSI Targets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_disk_pool_iscsi_target.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.StoragePool/diskPools/pool1/iscsiTargets/iscsiTarget1
```
