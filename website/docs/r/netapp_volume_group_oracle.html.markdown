---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume_group_oracle"
description: |-
  Manages a Application Volume Group for Oracle application.
---

# azurerm_netapp_volume_group_oracle

Manages a Application Volume Group for Oracle application.

~> **Note:** This feature is intended to be used for Oracle workloads only, with several requirements, please refer to [Understand Azure NetApp Files application volume group for Oracle](https://learn.microsoft.com/en-us/azure/azure-netapp-files/application-volume-oracle-introduction) document as the starting point to understand this feature before using it with Terraform.

## Example Usage

```hcl
provider "azurerm" {
  features {
    netapp {
      prevent_volume_destruction = true
    }
  }
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location

  tags = {
    "SkipNRMSNSG" = "true"
  }
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.88.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-delegated-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.88.2.0/24"]

  delegation {
    name = "exampledelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-netapp-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  depends_on = [
    azurerm_subnet.example
  ]
}

resource "azurerm_netapp_pool" "example" {
  name                = "${var.prefix}-netapp-pool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Manual"
}


resource "azurerm_netapp_volume_group_oracle" "example" {
  name                   = "${var.prefix}-NetAppVolumeGroupOracle"
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name
  account_name           = azurerm_netapp_account.example.name
  group_description      = "Example volume group for Oracle"
  application_identifier = "TST"

  volume {
    name                       = "${var.prefix}-volume-ora1"
    volume_path                = "${var.prefix}-my-unique-file-ora-path-1"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.example.id
    subnet_id                  = azurerm_subnet.example.id
    zone                       = "1"
    volume_spec_name           = "ora-data1"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }

  volume {
    name                       = "${var.prefix}-volume-oraLog"
    volume_path                = "${var.prefix}-my-unique-file-oralog-path"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.example.id
    subnet_id                  = azurerm_subnet.example.id
    zone                       = "1"
    volume_spec_name           = "ora-log"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `account_name` - (Required) Name of the account where the application volume group belong to. Changing this forces a new Application Volume Group to be created and data will be lost.

* `application_identifier` - (Required) The SAP System ID, maximum 3 characters, e.g. `OR1`. Changing this forces a new Application Volume Group to be created and data will be lost.

* `group_description` - (Required) Volume group description. Changing this forces a new Application Volume Group to be created and data will be lost.

* `location` - (Required) The Azure Region where the Application Volume Group should exist. Changing this forces a new Application Volume Group to be created and data will be lost.

* `name` - (Required) The name which should be used for this Application Volume Group. Changing this forces a new Application Volume Group to be created and data will be lost.

* `resource_group_name` - (Required) The name of the Resource Group where the Application Volume Group should exist. Changing this forces a new Application Volume Group to be created and data will be lost.

* `volume` - (Required) One or more `volume` blocks as defined below.

---

A `volume` block supports the following:

* `capacity_pool_id` - (Required) The ID of the Capacity Pool. Changing this forces a new Application Volume Group to be created and data will be lost.

* `name` - (Required) The name which should be used for this volume. Changing this forces a new Application Volume Group to be created and data will be lost.

* `protocols` - (Required) The target volume protocol expressed as a list. Changing this forces a new Application Volume Group to be created and data will be lost. Supported values for Application Volume Group include `NFSv3` or `NFSv4.1`.

* `proximity_placement_group_id` - (Optional) The ID of the proximity placement group (PPG). Changing this forces a new Application Volume Group to be created and data will be lost. 

~> **Note:** For Oracle application, it is required to have PPG enabled so Azure NetApp Files can pin the volumes next to your compute resources, please check [Requirements and considerations for application volume group for Oracle](https://learn.microsoft.com/en-us/azure/azure-netapp-files/application-volume-group-oracle-considerations) for details and other requirements. Note that this cannot be used together with `zone`.

* `zone` - (Optional) Specifies the Availability Zone in which the Volume should be located. Possible values are `1`, `2` and `3`, depending on the Azure region. Changing this forces a new resource to be created. This feature is currently in preview, for more information on how to enable it, please refer to [Manage availability zone volume placement for Azure NetApp Files](https://learn.microsoft.com/en-us/azure/azure-netapp-files/manage-availability-zone-volume-placement). Note that this cannot be used together with `proximity_placement_group_id`.

* `security_style` - (Required) Volume security style. Possible values are `ntfs` and `unix`. Changing this forces a new Application Volume Group to be created and data will be lost.

* `service_level` - (Required) Volume security style. Possible values are `Premium`, `Standard` and `Ultra`. Changing this forces a new Application Volume Group to be created and data will be lost.

* `snapshot_directory_visible` - (Required) Specifies whether the .snapshot (NFS clients) path of a volume is visible. Changing this forces a new Application Volume Group to be created and data will be lost.

* `storage_quota_in_gb` - (Required) The maximum Storage Quota allowed for a file system in Gigabytes.

* `subnet_id` - (Required) The ID of the Subnet the NetApp Volume resides in, which must have the `Microsoft.NetApp/volumes` delegation. Changing this forces a new Application Volume Group to be created and data will be lost.

* `throughput_in_mibps` - (Required) Throughput of this volume in Mibps.

* `volume_path` - (Required) A unique file path for the volume. Changing this forces a new Application Volume Group to be created and data will be lost.

* `volume_spec_name` - (Required) Volume specification name. Possible values are `ora-data1` through `ora-data8`, `ora-log`, `ora-log-mirror`, `ora-backup` and `ora-binary`. Changing this forces a new Application Volume Group to be created and data will be lost.

* `tags` - (Optional) A mapping of tags which should be assigned to the Application Volume Group.

* `network_features` - (Optional) Indicates which network feature to use, accepted values are `Basic` or `Standard`, it defaults to `Basic` if not defined. This is a feature in public preview and for more information about it and how to register, please refer to [Configure network features for an Azure NetApp Files volume](https://docs.microsoft.com/en-us/azure/azure-netapp-files/configure-network-features). This is required if enabling customer managed keys encryption scenario.

* `encryption_key_source` - (Optional) The encryption key source, it can be `Microsoft.NetApp` for platform managed keys or `Microsoft.KeyVault` for customer-managed keys. This is required with `key_vault_private_endpoint_id`. Changing this forces a new resource to be created.

* `key_vault_private_endpoint_id` - (Optional) The Private Endpoint ID for Key Vault, which is required when using customer-managed keys. This is required with `encryption_key_source`. Changing this forces a new resource to be created.

* `export_policy_rule` - (Required) One or more `export_policy_rule` blocks as defined below.

* `data_protection_snapshot_policy` - (Optional) A `data_protection_snapshot_policy` block as defined below.

---
A `data_protection_snapshot_policy` block supports the following:

* `snapshot_policy_id` - (Required) Resource ID of the snapshot policy to apply to the volume.

---

A `export_policy_rule` block supports the following:

* `allowed_clients` - (Required) A comma-sperated list of allowed client IPv4 addresses.

* `nfsv3_enabled` - (Required) Enables NFSv3. Please note that this cannot be enabled if volume has NFSv4.1 as its protocol.

* `nfsv41_enabled` - (Required) Enables NFSv4.1. Please note that this cannot be enabled if volume has NFSv3 as its protocol.

* `root_access_enabled` - (Optional) Is root access permitted to this volume? Defaults to `true`.

* `rule_index` - (Required) The index number of the rule, must start at 1 and maximum 5.

* `unix_read_only` - (Optional) Is the file system on unix read only? Defaults to `false.

* `unix_read_write` - (Optional) Is the file system on unix read and write? Defaults to `true`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Application Volume Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Application Volume Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Volume Group.
* `update` - (Defaults to 2 hours) Used when updating the Application Volume Group.
* `delete` - (Defaults to 2 hours) Used when deleting the Application Volume Group.

## Import

Application Volume Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_volume_group_oracle.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mytest-rg/providers/Microsoft.NetApp/netAppAccounts/netapp-account-test/volumeGroups/netapp-volumegroup-test
```
