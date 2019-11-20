---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume"
sidebar_current: "docs-azurerm-resource-netapp-volume"
description: |-
  Manages a NetApp Volume.
---

# azurerm_netapp_volume

Manages a NetApp Volume.


## NetApp Volume Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-virtualnetwork"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.2.0/24"
  delegation {
    name = "testdelegation"
  
    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "example" {
  name                = "example-netappaccount"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_netapp_pool" "example" {
  name                = "example-netapppool"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  account_name        = "${azurerm_netapp_account.example.name}"
  service_level       = "Premium"
  size_in_4_tb        = "1"
}

resource "azurerm_netapp_volume" "example" {
  name                = "example-netappvolume"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  account_name        = "${azurerm_netapp_account.example.name}"
  pool_name           = "${azurerm_netapp_pool.example.name}"
  creation_token      = "my-unique-file-path"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.example.id}"
  usage_threshold     = "100"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Volume. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Volume should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp account in which the NetApp Pool should be created.

* `pool_name` - (Required) The name of the NetApp pool in which the NetApp Volume should be created.

* `creation_token` - (Required) A unique file path for the volume. Used when creating mount targets.

* `service_level` - (Required) The service level of the file system. Valid values include `Premium`, `Standard`, or `Ultra`.

* `subnet_id` - (Required) The Azure Resource URI for a delegated subnet. Must have the delegation Microsoft.NetApp/volumes.

* `usage_threshold` - (Required) Maximum storage quota allowed for a file system in bytes. This is a soft quota used for alerting only.

* `export_policy_rule` - (Optional) One `export_policy_rule` block defined below.

---

The `export_policy_rule` block contains the following:

* `rule_index` - (Required) Order index.

* `allowed_clients` - (Required) Client ingress specification as comma separated string with IPv4 CIDRs, IPv4 host addresses and host names.

* `cifs` - (Required) Allows CIFS protocol.

* `nfsv3` - (Required) Allows NFSv3 protocol.

* `nfsv4` - (Required) Allows NFSv4 protocol.

* `unix_read_only` - (Required) Read only access.

* `unix_read_write` - (Required) Read and write access.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NetApp Volume.

## Import

NetApp Volume can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_volume.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1
```