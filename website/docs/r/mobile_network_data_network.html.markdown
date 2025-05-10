---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_data_network"
description: |-
  Manages a Mobile Network Data Network.
---

# azurerm_mobile_network_data_network

Manages a Mobile Network Data Network.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East Us"
}

resource "azurerm_mobile_network" "example" {
  name                = "example-mn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  mobile_country_code = "001"
  mobile_network_code = "01"
}

resource "azurerm_mobile_network_data_network" "example" {
  name              = "example-mndn"
  mobile_network_id = azurerm_mobile_network.example.id
  location          = azurerm_resource_group.example.location
  description       = "example description"

  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Data Network. Changing this forces a new Mobile Network Data Network to be created.

* `mobile_network_id` - (Required) Specifies the ID of the Mobile Network. Changing this forces a new Mobile Network Data Network to be created.

* `location` - (Required) Specifies the Azure Region where the Mobile Network Data Network should exist. Changing this forces a new Mobile Network Data Network to be created.

* `description` - (Optional) A description of this Mobile Network Data Network.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network Data Network.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Data Network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Mobile Network Data Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Data Network.
* `update` - (Defaults to 30 minutes) Used when updating the Mobile Network Data Network.
* `delete` - (Defaults to 3 hours) Used when deleting the Mobile Network Data Network.

## Import

Mobile Network Data Network can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_data_network.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/mobileNetworks/mobileNetwork1/dataNetworks/dataNetwork1
```
