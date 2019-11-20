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
  name                = "acctestnetappvolume"
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

* `creation_token` - A unique file path for the volume. Used when creating mount targets.

* `service_level` - The service level of the file system.

* `subnet_id` - The Azure Resource URI for a delegated subnet. Must have the delegation Microsoft.NetApp/volumes.

* `usage_threshold` - Maximum storage quota allowed for a file system in bytes. This is a soft quota used for alerting only.

* `export_policy_rule` - One `export_policy_rule` block defined below.

---

The `export_policy_rule` block contains the following:

* `rule_index` - Order index.

* `allowed_clients` - Client ingress specification as comma separated string with IPv4 CIDRs, IPv4 host addresses and host names.

* `cifs` - Allows CIFS protocol.

* `nfsv3` - Allows NFSv3 protocol.

* `nfsv4` - Allows NFSv4 protocol.

* `unix_read_only` - Read only access.

* `unix_read_write` - Read and write access.

---