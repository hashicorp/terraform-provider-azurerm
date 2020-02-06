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
  value = "${data.azurerm_netapp_pool.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Pool.

* `account_name` - (Required) The name of the NetApp account where the NetApp pool exists.

* `resource_group_name` - (Required) The Name of the Resource Group where the NetApp Pool exists.


## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the NetApp Pool exists.

* `service_level` - The service level of the file system.

* `size_in_tb` - Provisioned size of the pool in TB.

---

### Timeouts

~> **Note:** Custom Timeouts are available [as an opt-in Beta in version 1.43 of the Azure Provider](/docs/providers/azurerm/guides/2.0-beta.html) and will be enabled by default in version 2.0 of the Azure Provider.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Pool.
