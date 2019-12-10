---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_geographical_location"
sidebar_current: "docs-azurerm-datasource-traffic-manager-geographical-location"
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
