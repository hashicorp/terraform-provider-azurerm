---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maps_account"
sidebar_current: "docs-azurerm-resource-maps-account"
description: |-
  Manages an Azure Maps Account.
---

# azurerm_maps_account

Manages an Azure Maps Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_maps_account" "test" {
  name                = "example-maps-account"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku = "S1"

  tags = {
    environment = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Maps Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Azure Maps Account should exist. Changing this forces a new resource to be created.

* `sku` - (Optional) The sku of the Azure Maps Account. Values available are S0, S1. Defaults to S0.

* `tags` - (Optional) A mapping of tags to assign to the Azure Maps Account.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Azure Maps Account.

* `primary_access_key` - The primary key used to authenticate and authorize access to the Maps REST APIs.

* `secondary_access_key` - The primary key used to authenticate and authorize access to the Maps REST APIs. The second key is given to provide seemless key regeneration.

* `x_ms_client_id` - A unique identifier for the Maps Account.


## Import

API Management Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_maps_account.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Maps/accounts/my-maps-account
```

