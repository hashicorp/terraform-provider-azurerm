---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume"
description: |-
  Manages a NetApp Volume.
---

# azurerm_netapp_volume

Manages a NetApp Volume.

!> **Note:** This resource uses a feature to prevent deletion called `prevent_volume_destruction`, defaulting to `true`. It is intentionally set to `true` to prevent the possibility of accidental data loss. The example in this page shows all possible protection options you can apply, it is using same values as the defaults.

## NetApp Volume Usage

```hcl
provider "azurerm" {
  features {
    netapp {
      prevent_volume_destruction             = true
      delete_backups_on_backup_vault_destroy = false
    }
  }
}

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

resource "azurerm_netapp_backup_vault" "example" {
  name                = "example-netappbackupvault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  account_name        = azurerm_netapp_account.example.name
}

resource "azurerm_netapp_backup_policy" "example" {
  name                    = "example-netappbackuppolicy"
  resource_group_name     = azurerm_resource_group.example.name
  location                = azurerm_resource_group.example.location
  account_name            = azurerm_netapp_account.example.name
  daily_backups_to_keep   = 2
  weekly_backups_to_keep  = 2
  monthly_backups_to_keep = 2
  enabled                 = true
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
  name                       = "example-netappvolume"
  location                   = azurerm_resource_group.example.location
  zone                       = "1"
  resource_group_name        = azurerm_resource_group.example.name
  account_name               = azurerm_netapp_account.example.name
  pool_name                  = azurerm_netapp_pool.example.name
  volume_path                = "my-unique-file-path"
  service_level              = "Premium"
  subnet_id                  = azurerm_subnet.example.id
  protocols                  = ["NFSv4.1"]
  security_style             = "unix"
  storage_quota_in_gb        = 100
  snapshot_directory_visible = false

  # When creating volume from a snapshot
  create_from_snapshot_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1/snapshots/snapshot1"

  # Following section is only required if deploying a data protection volume (secondary)
  # to enable Cross-Region Replication feature
  data_protection_replication {
    endpoint_type             = "dst"
    remote_volume_location    = azurerm_resource_group.example.location
    remote_volume_resource_id = azurerm_netapp_volume.example.id
    replication_frequency     = "10minutes"
  }

  # Enabling Snapshot Policy for the volume
  # Note: this cannot be used in conjunction with data_protection_replication when endpoint_type is dst
  data_protection_snapshot_policy {
    snapshot_policy_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/snapshotPolicies/snapshotpolicy1"
  }

  # Enabling backup policy
  data_protection_backup_policy {
    backup_vault_id  = azurerm_netapp_backup_vault.example.id
    backup_policy_id = azurerm_netapp_backup_policy.example.id
    policy_enabled   = true
  }

  # prevent the possibility of accidental data loss
  lifecycle {
    prevent_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Volume. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Volume should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `zone` - (Optional) Specifies the Availability Zone in which the Volume should be located. Possible values are `1`, `2` and `3`. Changing this forces a new resource to be created. This feature is currently in preview, for more information on how to enable it, please refer to [Manage availability zone volume placement for Azure NetApp Files](https://learn.microsoft.com/en-us/azure/azure-netapp-files/manage-availability-zone-volume-placement#register-the-feature).

* `account_name` - (Required) The name of the NetApp account in which the NetApp Pool should be created. Changing this forces a new resource to be created.

* `volume_path` - (Required) A unique file path for the volume. Used when creating mount targets. Changing this forces a new resource to be created.

* `pool_name` - (Required) The name of the NetApp pool in which the NetApp Volume should be created.

* `service_level` - (Required) The target performance of the file system. Valid values include `Premium`, `Standard`, or `Ultra`.

~> **Note:** When updating `service_level` by migrating it to another Capacity Pool, both `service_level` and `pool_name` must be changed, otherwise the volume will be recreated with the specified `service_level`.

~> **Note:** After updating `service_level` the `id` for the volume will change to include the new Capacity Pool so any resources referencing the Volume will be silently removed from state. They will still exist in Azure but need to reimported into Terraform.

* `azure_vmware_data_store_enabled` - (Optional) Is the NetApp Volume enabled for Azure VMware Solution (AVS) datastore purpose. Defaults to `false`. Changing this forces a new resource to be created.

* `protocols` - (Optional) The target volume protocol expressed as a list. Supported single value include `CIFS`, `NFSv3`, or `NFSv4.1`. If argument is not defined it will default to `NFSv3`. Changing this forces a new resource to be created and data will be lost. Dual protocol scenario is supported for CIFS and NFSv3, for more information, please refer to [Create a dual-protocol volume for Azure NetApp Files](https://docs.microsoft.com/azure/azure-netapp-files/create-volumes-dual-protocol) document.

* `security_style` - (Optional) Volume security style, accepted values are `unix` or `ntfs`. If not provided, single-protocol volume is created defaulting to `unix` if it is `NFSv3` or `NFSv4.1` volume, if `CIFS`, it will default to `ntfs`. In a dual-protocol volume, if not provided, its value will be `ntfs`. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet the NetApp Volume resides in, which must have the `Microsoft.NetApp/volumes` delegation. Changing this forces a new resource to be created.

* `network_features` - (Optional) Indicates which network feature to use, accepted values are `Basic` or `Standard`, it defaults to `Basic` if not defined. This is a feature in public preview and for more information about it and how to register, please refer to [Configure network features for an Azure NetApp Files volume](https://docs.microsoft.com/en-us/azure/azure-netapp-files/configure-network-features).

* `storage_quota_in_gb` - (Required) The maximum Storage Quota allowed for a file system in Gigabytes.

* `snapshot_directory_visible` - (Optional) Specifies whether the .snapshot (NFS clients) or ~snapshot (SMB clients) path of a volume is visible, default value is true.

* `create_from_snapshot_resource_id` - (Optional) Creates volume from snapshot. Following properties must be the same as the original volume where the snapshot was taken from: `protocols`, `subnet_id`, `location`, `service_level`, `resource_group_name`, `account_name` and `pool_name`. Changing this forces a new resource to be created.

* `data_protection_replication` - (Optional) A `data_protection_replication` block as defined below. Changing this forces a new resource to be created.

* `data_protection_snapshot_policy` - (Optional) A `data_protection_snapshot_policy` block as defined below.

* `data_protection_backup_policy` - (Optional) A `data_protection_backup_policy` block as defined below.

* `export_policy_rule` - (Optional) One or more `export_policy_rule` block defined below.

* `throughput_in_mibps` - (Optional) Throughput of this volume in Mibps.

* `encryption_key_source` - (Optional) The encryption key source, it can be `Microsoft.NetApp` for platform managed keys or `Microsoft.KeyVault` for customer-managed keys. This is required with `key_vault_private_endpoint_id`. Changing this forces a new resource to be created.

* `kerberos_enabled` - (Optional) Enable to allow Kerberos secured volumes. Requires appropriate export rules.

~> **Note:** `kerberos_enabled` requires that the parent `azurerm_netapp_account` has a *valid* AD connection defined. If the configuration is invalid, the volume will still be created but in a failed state. This requires manually deleting the volume and recreating it again via Terraform once the AD configuration has been corrected.

* `key_vault_private_endpoint_id` - (Optional) The Private Endpoint ID for Key Vault, which is required when using customer-managed keys. This is required with `encryption_key_source`. Changing this forces a new resource to be created.

* `smb_non_browsable_enabled` - (Optional) Limits clients from browsing for an SMB share by hiding the share from view in Windows Explorer or when listing shares in "net view." Only end users that know the absolute paths to the share are able to find the share. Defaults to `false`. For more information, please refer to [Understand NAS share permissions in Azure NetApp Files](https://learn.microsoft.com/en-us/azure/azure-netapp-files/network-attached-storage-permissions#:~:text=Non%2Dbrowsable%20shares,find%20the%20share.)

* `smb_access_based_enumeration_enabled` - (Optional) Limits enumeration of files and folders (that is, listing the contents) in SMB only to users with allowed access on the share. For instance, if a user doesn't have access to read a file or folder in a share with access-based enumeration enabled, then the file or folder doesn't show up in directory listings. Defaults to `false`. For more information, please refer to [Understand NAS share permissions in Azure NetApp Files](https://learn.microsoft.com/en-us/azure/azure-netapp-files/network-attached-storage-permissions#:~:text=security%20for%20administrators.-,Access%2Dbased%20enumeration,in%20an%20Azure%20NetApp%20Files%20SMB%20volume.%20Only%20contosoadmin%20has%20access.,-In%20the%20below)

* `smb_continuous_availability_enabled` - (Optional) Enable SMB Continuous Availability.

* `smb3_protocol_encryption_enabled` - (Optional) Enable SMB encryption.

* `tags` - (Optional) A mapping of tags to assign to the resource.

-> **Note:** It is highly recommended to use the **lifecycle** property as noted in the example since it will prevent an accidental deletion of the volume if the `protocols` argument changes to a different protocol type.

---

An `export_policy_rule` block supports the following:

* `rule_index` - (Required) The index number of the rule.

* `allowed_clients` - (Required) A list of allowed clients IPv4 addresses.

* `protocols_enabled` - (Optional) A list of allowed protocols. Valid values include `CIFS`, `NFSv3`, or `NFSv4.1`. Only one value is supported at this time. This replaces the previous arguments: `cifs_enabled`, `nfsv3_enabled` and `nfsv4_enabled`.

* `unix_read_only` - (Optional) Is the file system on unix read only?

* `unix_read_write` - (Optional) Is the file system on unix read and write?

* `root_access_enabled` - (Optional) Is root access permitted to this volume?

* `kerberos_5_read_only_enabled` - (Optional) Is Kerberos 5 read-only access permitted to this volume?

* `kerberos_5_read_write_enabled` - (Optional) Is Kerberos 5 read/write permitted to this volume?

* `kerberos_5i_read_only_enabled` - (Optional) Is Kerberos 5i read-only permitted to this volume?

* `kerberos_5i_read_write_enabled` - (Optional) Is Kerberos 5i read/write permitted to this volume?

* `kerberos_5p_read_only_enabled` - (Optional) Is Kerberos 5p read-only permitted to this volume?

* `kerberos_5p_read_write_enabled` - (Optional) Is Kerberos 5p read/write permitted to this volume?

---

A `data_protection_replication` block is used when enabling the Cross-Region Replication (CRR) data protection option by deploying two Azure NetApp Files Volumes, one to be a primary volume and the other one will be the secondary, the secondary will have this block and will reference the primary volume, each volume must be in a supported [region pair](https://docs.microsoft.com/azure/azure-netapp-files/cross-region-replication-introduction#supported-region-pairs) and it supports the following:

* `endpoint_type` - (Optional) The endpoint type, default value is `dst` for destination.
  
* `remote_volume_location` - (Required) Location of the primary volume. Changing this forces a new resource to be created.

* `remote_volume_resource_id` - (Required) Resource ID of the primary volume.
  
* `replication_frequency` - (Required) Replication frequency, supported values are '10minutes', 'hourly', 'daily', values are case sensitive.

A full example of the `data_protection_replication` attribute can be found in [the `./examples/netapp/volume_crr` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/netapp/volume_crr)

~> **Note:** `data_protection_replication` can be defined only once per secondary volume, adding a second instance of it is not supported.

---

A `data_protection_snapshot_policy` block is used when automatic snapshots for a volume based on a specific snapshot policy. It supports the following:

* `snapshot_policy_id` - (Required) Resource ID of the snapshot policy to apply to the volume.

A full example of the `data_protection_snapshot_policy` attribute usage can be found in [the `./examples/netapp/nfsv3_volume_with_snapshot_policy` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/netapp/nfsv3_volume_with_snapshot_policy)
  
~> **Note:** `data_protection_snapshot_policy` block can be used alone or with data_protection_replication in the primary volume only, if enabling it in the secondary, an error will be thrown.

---

A `data_protection_backup_policy` block is used to setup automatic backups through a specific backup policy. It supports the following:

* `backup_vault_id` - (Required) Resource ID of the backup backup vault to associate this volume to.

* `backup_policy_id` - (Required) Resource ID of the backup policy to apply to the volume.

* `policy_enabled` - (Optional) Enables the backup policy on the volume, defaults to `true`.

For more information on Azure NetApp Files Backup feature please see [Understand Azure NetApp Files backup](https://learn.microsoft.com/en-us/azure/azure-netapp-files/backup-introduction)
  
---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NetApp Volume.

* `mount_ip_addresses` - A list of IPv4 Addresses which should be used to mount the volume.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the NetApp Volume.
* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Volume.
* `update` - (Defaults to 1 hour) Used when updating the NetApp Volume.
* `delete` - (Defaults to 1 hour) Used when deleting the NetApp Volume.

## Import

NetApp Volumes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_volume.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1
```
