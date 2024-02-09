---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_location"
description: |-
   Get information of a specific location.
---

# Data Source: azurerm_location

Use this data source to access information of a specific physical location.

## Example Usage

```hcl
data "azurerm_location" "example" {
  location = "West Europe"
}
```

## Argument Reference

* `location` - (Required) Specifies the supported Azure location where the resource exists.

## Attributes Reference

* `id` - The ID of Location within this Subscription.

* `region_type` - The available region type. Possible values are `Physical` and `Logical`.

* `display_name` - The display name of the location.

* `zone_mappings` - A `zone_mappings` block as defined below.

---

A `zone_mappings` block exports the following:

* `logical_zone` - The logical zone id for the availability zone

* `physical_zone` - The fully qualified physical zone id of availability zone to which logical zone id is mapped to


---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Location.
