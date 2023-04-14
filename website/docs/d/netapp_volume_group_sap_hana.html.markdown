---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_netapp_volume_group_sap_hana"
description: |-
  Gets information about an existing Application Volume Group for SAP HANA application.
---

# Data Source: azurerm_netapp_volume_group_sap_hana

Use this data source to access information about an existing Application Volume Group for SAP HANA application.

## Example Usage

```hcl
data "azurerm_netapp_volume_group_sap_hana" "example" {
  name                = "existing application volume group name"
  resource_group_name = "resource group name where the account and volume group belong to"
  account_name        = "existing account where the application volume group belong to"
}

output "id" {
  value = data.azurerm_netapp_volume_group_sap_hana.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `account_name` - (Required) Name of the account where the application volume group belong to.

* `name` - (Required) The name of this Application Volume Group for SAP HANA application.

* `resource_group_name` - (Required) The name of the Resource Group where the Application Volume Group exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Application Volume Group.

* `application_identifier` - The application identifier.

* `group_description` - Volume group description.

* `location` - The Azure Region where the Application Volume Group exists.

* `volume` - A `volume` block as defined below.

---

A `volume` block exports the following:

* `capacity_pool_id` - The ID of the Capacity Pool.

* `id` - Volume ID.

* `name` - The name of this volume.

* `proximity_placement_group_id` - The ID of the proximity placement group.

* `security_style` - Volume security style.

* `service_level` - The target performance of the file system.

* `snapshot_directory_visible` - Is the .snapshot (NFS clients) path of a volume visible?

* `storage_quota_in_gb` - The maximum Storage Quota allowed for a file system in Gigabytes.

* `subnet_id` - The ID of the Subnet the NetApp Volume resides in.

* `tags` - A mapping of tags assigned to the Application Volume Group.

* `throughput_in_mibps` - Throughput of this volume in Mibps.

* `volume_path` - A unique file path for the volume.

* `volume_spec_name` - Volume spec name.

* `data_protection_replication` - A `data_protection_replication` block as defined below.

* `data_protection_snapshot_policy` - A `data_protection_snapshot_policy` block as defined below.

* `export_policy_rule` - A `export_policy_rule` block as defined below.

* `mount_ip_addresses` - A `mount_ip_addresses` block as defined below.

* `protocols` - A `protocols` block as defined below.

---

A `data_protection_replication` block exports the following:

* `endpoint_type` - The endpoint type.

* `remote_volume_location` - Location of the primary volume.

* `remote_volume_resource_id` - Resource ID of the primary volume.

* `replication_frequency` - Replication frequency.

---

A `data_protection_snapshot_policy` block exports the following:

* `snapshot_policy_id` - Resource ID of the snapshot policy to apply to the volume.

---

A `export_policy_rule` block exports the following:

* `allowed_clients` - A list of allowed clients IPv4 addresses.

* `nfsv3_enabled` - Is the NFSv3 protocol enabled?

* `nfsv41_enabled` - Is the NFSv4.1 enabled?

* `root_access_enabled` - Is root access permitted to this volume?

* `rule_index` - The index number of the rule.

* `unix_read_only` - Is the file system on unix read only?.

* `unix_read_write` - Is the file system on unix read and write?.

---



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Application Volume Group.
