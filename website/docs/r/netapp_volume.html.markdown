---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume"
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
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
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
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Premium"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "example" {
  lifecycle {
    prevent_destroy = true
  }

  name                = "example-netappvolume"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "my-unique-file-path"
  service_level       = "Premium"
  subnet_id           = azurerm_subnet.example.id
  protocols           = ["NFSv4.1"]
  storage_quota_in_gb = 100

  data_protection_replication {
    endpoint_type             = "dst"
    remote_volume_location    = azurerm_resource_group.example_primary.location
    remote_volume_resource_id = azurerm_netapp_volume.example_primary.id
    replication_schedule      = "_10minutely"
  }
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

* `protocols` - (Optional) The target volume protocol expressed as a list. Supported single value include `CIFS`, `NFSv3`, or `NFSv4.1`. If argument is not defined it will default to `NFSv3`. Changing this forces a new resource to be created and data will be lost.

* `subnet_id` - (Required) The ID of the Subnet the NetApp Volume resides in, which must have the `Microsoft.NetApp/volumes` delegation. Changing this forces a new resource to be created.

* `storage_quota_in_gb` - (Required) The maximum Storage Quota allowed for a file system in Gigabytes.

* `export_policy_rule` - (Optional) One or more `export_policy_rule` block defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

-> **Note**: It is highly recommended to use the **lifecycle** property as noted in the example since it will prevent an accidental deletion of the volume if the `protocols` argument changes to a different protocol type. 

---

An `export_policy_rule` block supports the following:

* `rule_index` - (Required) The index number of the rule.

* `allowed_clients` - (Required) A list of allowed clients IPv4 addresses.

* `protocols_enabled` - (Optional) A list of allowed protocols. Valid values include `CIFS`, `NFSv3`, or `NFSv4.1`. Only one value is supported at this time. This replaces the previous arguments: `cifs_enabled`, `nfsv3_enabled` and `nfsv4_enabled`.

* `cifs_enabled` - (Optional / **Deprecated in favour of `protocols_enabled`**) Is the CIFS protocol allowed?

* `nfsv3_enabled` - (Optional / **Deprecated in favour of `protocols_enabled`**) Is the NFSv3 protocol allowed?

* `nfsv4_enabled` - (Optional / **Deprecated in favour of `protocols_enabled`**)  Is the NFSv4 protocol allowed?

* `unix_read_only` - (Optional) Is the file system on unix read only?

* `unix_read_write` - (Optional) Is the file system on unix read and write?

---

An `data_protection_replication` is used when enabling the Cross-Region Replication (CRR) data protection option by deploying two Azure NetApp Files Volumes, one to be a primary volume and the other one will be the secondary, the secondary will have this block and will reference the primary volume, each volume must be in a supported region pair and it supports the following:

* `endpoint_type` - (Required) The endpoint type, supported value is `dst` for destination.
  
* `remote_volume_location` - (Required) Primary volume's location.

* `remote_volume_resource_id` - (Required) Primary volume's resource id.
  
* `replication_schedule` - (Required) Replication frequency, supported values are '_10minutely', 'hourly', 'daily'

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NetApp Volume.

* `mount_ip_addresses` - A list of IPv4 Addresses which should be used to mount the volume.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the NetApp Volume.
* `update` - (Defaults to 30 minutes) Used when updating the NetApp Volume.
* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Volume.
* `delete` - (Defaults to 30 minutes) Used when deleting the NetApp Volume.

## Import

NetApp Volumes can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_volume.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1
```
