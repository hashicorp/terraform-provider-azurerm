---
subcategory: "Maps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maps_account"
description: |-
  Manages an Azure Maps Account.
---

# azurerm_maps_account

Manages an Azure Maps Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_maps_account" "example" {
  name                = "example-maps-account"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "S1"

  tags = {
    environment = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Maps Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Azure Maps Account should exist. Changing this forces a new resource to be created.

* `sku_name` - (Required) The sku of the Azure Maps Account. Possible values are `S0` and `S1`.

* `tags` - (Optional) A mapping of tags to assign to the Azure Maps Account.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Azure Maps Account.

* `primary_access_key` - The primary key used to authenticate and authorize access to the Maps REST APIs.

* `secondary_access_key` - The secondary key used to authenticate and authorize access to the Maps REST APIs.

* `x_ms_client_id` - A unique identifier for the Maps Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Maps Account.
* `update` - (Defaults to 30 minutes) Used when updating the Maps Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Maps Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Maps Account.

## Import

A Maps Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_maps_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Maps/accounts/my-maps-account
```
