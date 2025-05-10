---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_site"
description: |-
  Manages a Mobile Network Site.
---

# azurerm_mobile_network_site

Manages a Mobile Network Site.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_databox_edge_device" "example" {
  name                = "example-device"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "EdgeP_Base-Standard"
}

resource "azurerm_mobile_network" "example" {
  name                = "example-mn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  mobile_country_code = "001"
  mobile_network_code = "01"
}

resource "azurerm_mobile_network_site" "example" {
  name              = "example-mns"
  mobile_network_id = azurerm_mobile_network.example.id
  location          = azurerm_resource_group.example.location

  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Mobile Network Site. Changing this forces a new Mobile Network Site to be created.

* `mobile_network_id` - (Required) the ID of the Mobile Network which the Mobile Network Site belongs to. Changing this forces a new Mobile Network Site to be created.

* `location` - (Required) The Azure Region where the Mobile Network Site should exist. Changing this forces a new Mobile Network Site to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network Site.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Site.

* `network_function_ids` - An array of Id of Network Functions deployed on the site.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Mobile Network Site.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Site.
* `update` - (Defaults to 3 hours) Used when updating the Mobile Network Site.
* `delete` - (Defaults to 3 hours) Used when deleting the Mobile Network Site.

## Import

Mobile Network Site can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_site.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/mobileNetworks/mobileNetwork1/sites/site1
```
