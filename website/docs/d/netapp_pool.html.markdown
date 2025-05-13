---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_pool"
description: |-
  Gets information about an existing NetApp Pool
---

# Data Source: azurerm_netapp_pool

Uses this data source to access information about an existing NetApp Pool.

## NetApp Pool Usage

```hcl
data "azurerm_netapp_pool" "example" {
  resource_group_name = "acctestRG"
  account_name        = "acctestnetappaccount"
  name                = "acctestnetapppool"
}

output "netapp_pool_id" {
  value = data.azurerm_netapp_pool.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the NetApp Pool.

* `account_name` - The name of the NetApp account where the NetApp pool exists.

* `resource_group_name` - The Name of the Resource Group where the NetApp Pool exists.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the NetApp Pool exists.

* `service_level` - The service level of the file system.

* `size_in_tb` - Provisioned size of the pool in TB.

* `encryption_type` - The encryption type of the pool.

* `cool_access_enabled` - Whether the NetApp Pool can hold cool access enabled volumes.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Pool.
