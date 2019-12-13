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
    name = "netapp"
  
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
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "example" {
  name                = "example-netappvolume"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  account_name        = "${azurerm_netapp_account.example.name}"
  pool_name           = "${azurerm_netapp_pool.example.name}"
  volume_path         = "my-unique-file-path"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.example.id}"
  storage_quota_in_gb = 100
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Volume. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Volume should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp account in which the NetApp Pool should be created. Changing this forces a new resource to be created.

* `volume_path` - (Required) A unique file path for the volume. Used when creating mount targets. Changing this forces a new resource to be created.

* `pool_name` - (Required) The name of the NetApp pool in which the NetApp Volume should be created. Changing this forces a new resource to be created.

* `service_level` - (Required) The target performance of the file system. Valid values include `Premium`, `Standard`, or `Ultra`.

* `subnet_id` - (Required) The ID of the Subnet the NetApp Volume resides in, which must have the `Microsoft.NetApp/volumes` delegation. Changing this forces a new resource to be created.

* `storage_quota_in_gb` - (Required) The maximum Storage Quota allowed for a file system in Gigabytes.

* `export_policy_rule` - (Optional) One or more `export_policy_rule` block defined below.

---

An `export_policy_rule` block supports the following:

* `rule_index` - (Required) The index number of the rule.

* `allowed_clients` - (Required) A list of allowed clients IPv4 addresses.

* `cifs_enabled` - (Required) Is the CIFS protocol allowed?

* `nfsv3_enabled` - (Required) Is the NFSv3 protocol allowed?

* `nfsv4_enabled` - (Required) Is the NFSv4 protocol allowed?

* `unix_read_only` - (Optional) Is the file system on unix read only?

* `unix_read_write` - (Optional) Is the file system on unix read and write?

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NetApp Volume.

## Import

NetApp Volumes can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_volume.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1
```