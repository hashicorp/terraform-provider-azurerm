---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_pool"
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
  size_in_tb          = 4
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Pool. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Pool should be created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp account in which the NetApp Pool should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `service_level` - (Required) The service level of the file system. Valid values include `Premium`, `Standard`, and `Ultra`. Changing this forces a new resource to be created.

* `size_in_tb` - (Required) Provisioned size of the pool in TB. Value must be between `1` and `2048`.

~> **Note:** `2` TB capacity pool sizing is currently in preview. You can only take advantage of the `2` TB minimum if all the volumes in the capacity pool are using `Standard` network features. If any volume is using `Basic` network features, the minimum size is `4` TB. Please see the product [documentation](https://learn.microsoft.com/azure/azure-netapp-files/azure-netapp-files-set-up-capacity-pool) for more information.

~> **Note:** The maximum `size_in_tb` is goverened by regional quotas. You may request additional capacity from Azure, currently up to `2048`.

* `qos_type` - (Optional) QoS Type of the pool. Valid values include `Auto` or `Manual`. Defaults to `Auto`.

* `encryption_type` - (Optional) The encryption type of the pool. Valid values include `Single`, and `Double`. Defaults to `Single`. Changing this forces a new resource to be created.

* `cool_access_enabled` - (Optional) Whether the NetApp Pool can hold cool access enabled volumes. Defaults to `false`.

~> **Note:** Disabling `cool_access_enabled` is not allowed and forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NetApp Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the NetApp Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Pool.
* `update` - (Defaults to 30 minutes) Used when updating the NetApp Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the NetApp Pool.

## Import

NetApp Pool can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1
```
