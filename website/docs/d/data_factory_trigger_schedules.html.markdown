---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_factory_trigger_schedules"
description: |-
  Gets information about all existing trigger schedules in Azure Data Factory.
---

# Data Source: azurerm_data_factory_trigger_schedules

Use this data source to access information about all existing trigger schedules in Azure Data Factory.

## Example Usage

```hcl
data "azurerm_data_factory_trigger_schedules" "example" {
  data_factory_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DataFactory/factories/datafactory1"
}

output "items" {
  value = data.azurerm_data_factory_trigger_schedules.example.items
}
```

## Arguments Reference

The following arguments are supported:

- `data_factory_id` - (Required) The ID of the Azure Data Factory to fetch trigger schedules from.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Azure Data Factory.

- `items` - A list of trigger schedule names available in this Azure Data Factory.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Data Factory trigger schedules.
