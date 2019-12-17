---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_platform_image"
sidebar_current: "docs-azurerm-datasource-platform-image"
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
  value = "${data.azurerm_platform_image.example.version}"
}
```

## Argument Reference

* `location` - (Required) Specifies the Location to pull information about this Platform Image from.
* `publisher` - (Required) Specifies the Publisher associated with the Platform Image.
* `offer` - (Required) Specifies the Offer associated with the Platform Image.
* `sku` - (Required) Specifies the SKU of the Platform Image.


## Attributes Reference

* `id` - The ID of the Platform Image.
* `version` - The latest version of the Platform Image.
