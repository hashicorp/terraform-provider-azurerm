---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume"
sidebar_current: "docs-azurerm-datasource-netapp-volume"
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
  value = "${data.azurerm_netapp_volume.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Volume.

* `resource_group_name` - (Required) The Name of the Resource Group where the NetApp Volume exists.

* `account_name` - (Required) The name of the NetApp account where the NetApp pool exists.

* `pool_name` - (Required) The name of the NetApp pool where the NetApp volume exists.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the NetApp Volume exists.

* `volume_path` - The unique file path of the volume.

* `service_level` - The service level of the file system.

* `subnet_id` - The ID of a Subnet in which the NetApp Volume resides.

* `storage_quota_in_gb` - The maximum Storage Quota in Gigabytes allowed for a file system.
