---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume"
description: |-
  Gets information about an existing NetApp Volume
---

# Data Source: azurerm_netapp_volume

Uses this data source to access information about an existing NetApp Volume.

## NetApp Volume Usage

```hcl
data "azurerm_netapp_volume" "example" {
  resource_group_name = "acctestRG"
  account_name        = "acctestnetappaccount"
  pool_name           = "acctestnetapppool"
  name                = "example-volume"
}

output "netapp_volume_id" {
  value = data.azurerm_netapp_volume.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the NetApp Volume.

* `resource_group_name` - The Name of the Resource Group where the NetApp Volume exists.

* `account_name` - The name of the NetApp account where the NetApp pool exists.

* `pool_name` - The name of the NetApp pool where the NetApp volume exists.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the NetApp Volume exists.

* `mount_ip_addresses` - A list of IPv4 Addresses which should be used to mount the volume.

* `protocols` - A list of protocol types enabled on volume.

* `service_level` - The service level of the file system.

* `subnet_id` - The ID of a Subnet in which the NetApp Volume resides.

* `storage_quota_in_gb` - The maximum Storage Quota in Gigabytes allowed for a file system.
 
* `security_style` - Volume security style

* `data_protection_replication` - Volume data protection block
* 
* `volume_path` - The unique file path of the volume.

---

A `data_protection_replication` block exports the following:

* `endpoint_type` - The endpoint type.

* `remote_volume_location` - Location of the primary volume.

* `remote_volume_resource_id` - Resource ID of the primary volume.

* `replication_frequency` - Frequency of replication.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Volume.
