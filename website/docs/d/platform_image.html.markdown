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

output "id" {
  value = data.azurerm_platform_image.example.id
}
```

## Argument Reference

* `location` - (Required) Specifies the Location to pull information about this Platform Image from.

* `publisher` - (Required) Specifies the Publisher associated with the Platform Image.

* `offer` - (Required) Specifies the Offer associated with the Platform Image.

* `sku` - (Required) Specifies the SKU of the Platform Image.

* `version` - (Optional) The version of the Platform Image.

## Attributes Reference

* `id` - The ID of the Platform Image.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Platform Image.
