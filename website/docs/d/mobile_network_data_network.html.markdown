---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_data_network"
description: |-
  Get information about a Mobile Network Data Network.
---

# azurerm_mobile_network_data_network

Get information about a Mobile Network Data Network.

## Example Usage

```hcl
data "azurerm_mobile_network" "example" {
  name                = "example-mn"
  resource_group_name = "example-rg"
}

data "azurerm_mobile_network_data_network" "example" {
  name              = "example-mndn"
  mobile_network_id = data.azurerm_mobile_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Data Network. 

* `mobile_network_id` - (Required) Specifies the ID of the Mobile Network. 

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Data Network.

* `location` - The Azure Region where the Mobile Network Data Network exists. 

* `description` - The description for this Mobile Network Data Network.

* `tags` - A mapping of tags which should be assigned to the Mobile Network Data Network.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Data Network.
