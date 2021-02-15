---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_snapshot"
description: |-
  Manages a NetApp Snapshot.
---

# azurerm_netapp_snapshot

Manages a NetApp Snapshot.

## NetApp Snapshot Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-virtualnetwork"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "netapp"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "example" {
  name                = "example-netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_pool" "example" {
  name                = "example-netapppool"
  account_name        = azurerm_netapp_account.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_level       = "Premium"
  size_in_tb          = "4"
}

resource "azurerm_netapp_volume" "example" {
  name                = "example-netappvolume"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "my-unique-file-path"
  service_level       = "Premium"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = "100"
}

resource "azurerm_netapp_snapshot" "example" {
  name                = "example-netappsnapshot"
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_name         = azurerm_netapp_volume.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Snapshot. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Snapshot should be created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp account in which the NetApp Pool should be created. Changing this forces a new resource to be created.

* `pool_name` - (Required) The name of the NetApp pool in which the NetApp Volume should be created. Changing this forces a new resource to be created.

* `volume_name` - (Required) The name of the NetApp volume in which the NetApp Snapshot should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NetApp Snapshot.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the NetApp Snapshot.
* `update` - (Defaults to 30 minutes) Used when updating the NetApp Snapshot.
* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Snapshot.
* `delete` - (Defaults to 30 minutes) Used when deleting the NetApp Snapshot.

## Import

NetApp Snapshot can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_snapshot.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1/snapshots/snapshot1
```
