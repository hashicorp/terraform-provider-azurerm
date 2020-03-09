---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
description: |-
  Gets information about an existing Virtual Hub
---

# Data Source: azurerm_virtual_hub

Uses this data source to access information about an existing Virtual Hub.

## Virtual Hub Usage

```hcl
data "azurerm_virtual_hub" "example" {
  name                = "example-hub"
  resource_group_name = "example-resources"
}

output "virtual_hub_id" {
  value = data.azurerm_virtual_hub.example.id
}
```


## Argument Reference

The following arguments are supported:

* `name` - The name of the Virtual Hub.

* `resource_group_name` - The Name of the Resource Group where the Virtual Hub exists.


## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the Virtual Hub exists.

* `address_prefix` - The Address Prefix used for this Virtual Hub.

* `tags` - A mapping of tags assigned to the Virtual Hub.

* `virtual_wan_id` - The ID of the Virtual WAN within which the Virtual Hub exists.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Hub.
