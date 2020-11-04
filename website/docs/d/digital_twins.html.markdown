---
subcategory: "DigitalTwins"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_digital_twins"
description: |-
  Gets information about an existing Digital Twins.
---

# Data Source: azurerm_digital_twins

Use this data source to access information about an existing Digital Twins.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_digital_twins" "example" {
  name                = "existing-digital-twins"
  resource_group_name = "existing-resgroup"
}

output "id" {
  value = data.azurerm_digital_twins.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Digital Twins.

* `resource_group_name` - (Required) The name of the Resource Group where the Digital Twins exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Digital Twins.

* `host_name` - The Api endpoint to work with this Digital Twins.

* `location` - The Azure Region where the Digital Twins exists.

* `tags` - A mapping of tags assigned to the Digital Twins.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Digital Twins.
