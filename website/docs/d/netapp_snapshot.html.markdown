---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_snapshot"
description: |-
  Gets information about an existing NetApp Snapshot
---

# Data Source: azurerm_netapp_snapshot

Uses this data source to access information about an existing NetApp Snapshot.

## NetApp Snapshot Usage

```hcl
data "azurerm_netapp_snapshot" "test" {
  resource_group_name = "acctestRG"
  name                = "acctestnetappsnapshot"
  account_name        = "acctestnetappaccount"
  pool_name           = "acctestnetapppool"
  volume_name         = "acctestnetappvolume"
}

output "netapp_snapshot_id" {
  value = data.azurerm_netapp_snapshot.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the NetApp Snapshot.

* `account_name` - The name of the NetApp Account where the NetApp Pool exists.

* `pool_name` - The name of the NetApp Pool where the NetApp Volume exists.

* `volume_name` - The name of the NetApp Volume where the NetApp Snapshot exists.

* `resource_group_name` - The Name of the Resource Group where the NetApp Snapshot exists.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the NetApp Snapshot exists.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Snapshot.
