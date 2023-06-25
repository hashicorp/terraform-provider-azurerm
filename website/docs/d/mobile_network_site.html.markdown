---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_site"
description: |-
  Get information about a Mobile Network Site.
---

# azurerm_mobile_network_site

Get information about a Mobile Network Site.

## Example Usage

```hcl
data "azurerm_mobile_network" "example" {
  name                = "example-mn"
  resource_group_name = "example-rg"
}

data "azurerm_mobile_network_site" "example" {
  name              = "example-mns"
  mobile_network_id = data.azurerm_mobile_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Mobile Network Site. 

* `mobile_network_id` - (Required) the ID of the Mobile Network which the Mobile Network Site belongs to. 

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Site.

* `network_function_ids` - An array of Id of Network Functions deployed on the site.

* `location` - The Azure Region where the Mobile Network Site should exist.

* `tags` - A mapping of tags which should be assigned to the Mobile Network Site.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Site.
