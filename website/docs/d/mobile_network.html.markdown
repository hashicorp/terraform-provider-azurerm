---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network"
description: |-
  Get information about an Azure Mobile Network.
---

# azurerm_mobile_network

Get information about an Azure Mobile Network.

## Example Usage

```hcl

data "azurerm_mobile_network" "example" {
  name                = "example-mn"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - Specifies the name which should be used for this Mobile Network.

* `resource_group_name` - Specifies the name of the Resource Group where the Mobile Network should exist. 

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network.

* `location` - The Azure Region where the Mobile Network should exist. Changing this forces a new Mobile Network to be created.

* `mobile_country_code` - Mobile country code (MCC), defined in https://www.itu.int/rec/T-REC-E.212 .

* `mobile_network_code` - Mobile network code (MNC), defined in https://www.itu.int/rec/T-REC-E.212 .

* `tags` - A mapping of tags which should be assigned to the Mobile Network.

* `service_key` - The mobile network resource identifier.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network.

## Import

Mobile Network can be imported using the `resource id`, e.g.
