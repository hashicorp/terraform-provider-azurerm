---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_pool"
sidebar_current: "docs-azurerm-resource-netapp-pool"
description: |-
  Manages a Pool within a NetApp Account.
---

# azurerm_netapp_pool

Manages a Pool within a NetApp Account.

## NetApp Pool Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_netapp_account" "example" {
  name                = "example-netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_pool" "example" {
  name                = "example-netapppool"
  account_name        = azurerm_netapp_account.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_level       = "Premium"
  size_in_tb          = "4"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Pool. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Pool should be created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp account in which the NetApp Pool should be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `service_level` - (Required) The service level of the file system. Valid values include `Premium`, `Standard`, or `Ultra`.

* `size_in_tb` - (Required) Provisioned size of the pool in TB. Value must be between `4` and `500`.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NetApp Pool.

## Import

NetApp Pool can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1
```