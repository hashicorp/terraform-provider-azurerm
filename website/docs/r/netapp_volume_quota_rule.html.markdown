---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume_quota_rule"
description: |-
  Manages a Volume Quota Rule.
---

# azurerm_netapp_volume_quota_rule

Manages a Volume Quota Rule.

## Example Usage

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

  name                       = "example-netappvolume"
  location                   = azurerm_resource_group.example.location
  zone                       = "1"
  resource_group_name        = azurerm_resource_group.example.name
  account_name               = azurerm_netapp_account.example.name
  pool_name                  = azurerm_netapp_pool.example.name
  volume_path                = "my-unique-file-path"
  service_level              = "Premium"
  subnet_id                  = azurerm_subnet.example.id
  network_features           = "Basic"
  protocols                  = ["NFSv4.1"]
  security_style             = "unix"
  storage_quota_in_gb        = 100
  snapshot_directory_visible = false
}

resource "azurerm_netapp_volume_quota_rule" "quota1" {
  name              = "example-quota-rule-1"
  location          = azurerm_resource_group.example.location
  volume_id         = azurerm_netapp_volume.example.id
  quota_target      = "3001"
  quota_size_in_kib = 1024
  quota_type        = "IndividualGroupQuota"
}

resource "azurerm_netapp_volume_quota_rule" "quota2" {
  name              = "example-quota-rule-2"
  location          = azurerm_resource_group.example.location
  volume_id         = azurerm_netapp_volume.example.id
  quota_target      = "2001"
  quota_size_in_kib = 1024
  quota_type        = "IndividualUserQuota"
}

resource "azurerm_netapp_volume_quota_rule" "quota3" {
  name              = "example-quota-rule-3"
  location          = azurerm_resource_group.example.location
  volume_id         = azurerm_netapp_volume.example.id
  quota_size_in_kib = 1024
  quota_type        = "DefaultUserQuota"
}

resource "azurerm_netapp_volume_quota_rule" "quota4" {
  name              = "example-quota-rule-4"
  location          = azurerm_resource_group.example.location
  volume_id         = azurerm_netapp_volume.example.id
  quota_size_in_kib = 1024
  quota_type        = "DefaultGroupQuota"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Volume Quota Rule should exist. Changing this forces a new Volume Quota Rule to be created.

* `name` - (Required) The name which should be used for this Volume Quota Rule. Changing this forces a new Volume Quota Rule to be created.

* `volume_id` - (Required) The NetApp volume ID where the Volume Quota Rule is assigned to. Changing this forces a new resource to be created.

* `quota_size_in_kib` - (Required) Quota size in kibibytes.

* `quota_type` - (Required) Quota type. Possible values are `DefaultGroupQuota`, `DefaultUserQuota`, `IndividualGroupQuota` and `IndividualUserQuota`. Please note that `IndividualGroupQuota` and `DefaultGroupQuota` are not applicable to SMB and dual-protocol volumes. Changing this forces a new resource to be created.

* `quota_target` - (Optional) Quota Target. This can be Unix UID/GID for NFSv3/NFSv4.1 volumes and Windows User SID for CIFS based volumes. Changing this forces a new resource to be created.

-> **Note:** `quota_target ` must be used when `quota_type` is `IndividualGroupQuota` or `IndividualUserQuota`

~> **Note:** more information about this resource can be found at [Understand default and individual user and group quotas](https://learn.microsoft.com/en-us/azure/azure-netapp-files/default-individual-user-group-quotas-introduction)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Volume Quota Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Volume Quota Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Volume Quota Rule.
* `update` - (Defaults to 2 hours) Used when updating the Volume Quota Rule.
* `delete` - (Defaults to 2 hours) Used when deleting the Volume Quota Rule.

## Import

Volume Quota Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_volume_quota_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/vol1/volumeQuotaRules/quota1
```
