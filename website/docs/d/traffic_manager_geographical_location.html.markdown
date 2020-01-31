---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_geographical_location"
description: |-
  Gets information about a specified Traffic Manager Geographical Location within the Geographical Hierarchy.

---

# Data Source: azurerm_traffic_manager_geographical_location

Use this data source to access the ID of a specified Traffic Manager Geographical Location within the Geographical Hierarchy.

## Example Usage (World)

```hcl
data "azurerm_traffic_manager_geographical_location" "example" {
  name = "World"
}

output "location_code" {
  value = "${data.azurerm_traffic_manager_geographical_location.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Location, for example `World`, `Europe` or `Germany`.

## Attributes Reference

* `id` - The ID of this Location, also known as the `Code` of this Location.

### Timeouts

~> **Note:** Custom Timeouts are available [as an opt-in Beta in version 1.43 of the Azure Provider](/docs/providers/azurerm/guides/2.0-beta.html) and will be enabled by default in version 2.0 of the Azure Provider.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the   Gets information about a specified Traffic Manager Geographical Location within the Geographical Hierarchy.
