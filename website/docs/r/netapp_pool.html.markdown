---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_pool"
sidebar_current: "docs-azurerm-resource-netapp-pool"
description: |-
  Manages a NetApp Pool.
---

# azurerm_netapp_pool

Manages a NetApp Pool.


## NetApp Pool Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_netapp_account" "example" {
  name                = "example-netappaccount"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_netapp_pool" "example" {
  name                = "example-netapppool"
  account_name        = "${azurerm_netapp_account.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  service_level       = "Premium"
  size                = "4398046511104"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Pool. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Pool should be created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp account.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `service_level` - (Required) The service level of the file system.

* `size` - (Required) Provisioned size of the pool (in bytes). Allowed values are in 4TiB chunks (value must be multiply of 4398046511104).

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NetApp Pool.

## Import

NetApp Pool can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1
```