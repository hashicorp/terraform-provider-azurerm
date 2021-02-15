---
subcategory: "Digital Twins"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_digital_twins_instance"
description: |-
  Gets information about an existing Digital Twins instance.
---

# Data Source: azurerm_digital_twins_instance

Use this data source to access information about an existing Digital Twins instance.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_digital_twins_instance" "example" {
  name                = "existing-digital-twins"
  resource_group_name = "existing-resgroup"
}

output "id" {
  value = data.azurerm_digital_twins_instance.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Digital Twins instance.

* `resource_group_name` - (Required) The name of the Resource Group where the Digital Twins instance exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Digital Twins instance.

* `host_name` - The Api endpoint to work with this Digital Twins instance.

* `location` - The Azure Region where the Digital Twins instance exists.

* `tags` - A mapping of tags assigned to the Digital Twins instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Digital Twins instance.
