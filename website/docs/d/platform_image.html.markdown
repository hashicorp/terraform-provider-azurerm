---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_platform_image"
description: |-
  Gets information about a Platform Image.
---

# Data Source: azurerm_platform_image

Use this data source to access information about a Platform Image.

## Example Usage

```hcl
data "azurerm_platform_image" "example" {
  location  = "West Europe"
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "16.04-LTS"
}

output "version" {
  value = data.azurerm_platform_image.example.version
}
```

## Argument Reference

* `location` - Specifies the Location to pull information about this Platform Image from.
* `publisher` - Specifies the Publisher associated with the Platform Image.
* `offer` - Specifies the Offer associated with the Platform Image.
* `sku` - Specifies the SKU of the Platform Image.


## Attributes Reference

* `id` - The ID of the Platform Image.
* `version` - The latest version of the Platform Image.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Platform Image.
