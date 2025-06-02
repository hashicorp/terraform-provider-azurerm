---
subcategory: "Extended Location"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_extended_location_custom_location"
description: |-
  Gets information about an existing Custom Location within an Extended Location.
---

# Data Source: azurerm_extended_location_custom_location

Use this data source to access information about an existing Custom Location within an Extended Location.

## Example Usage

```hcl
data "azurerm_extended_location_custom_location" "example" {
  name                = azurerm_extended_location_custom_location.example.name
  resource_group_name = azurerm_resource_group.example.name
}

output "custom_location_id" {
  value = data.azurerm_extended_location_custom_location.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Custom Location.

* `resource_group_name` - (Required) The name of the Resource Group where the Custom Location exists.

## Attributes Reference

* `id` - The ID of the Custom Location.

* `location` - The Azure location where the Custom Location exists.

* `namespace` - The namespace of the Custom Location.

* `cluster_extension_ids` - The list of Cluster Extension IDs.

* `host_resource_id` - The host resource ID.

* `authentication` - An `authentication` block as defined below.

* `display_name` - The display name of the Custom Location.

* `host_type` - The host type of the Custom Location.

---

An `authentication` block supports the following:

* `type` - The type of authentication.

* `value` - The value of authentication.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Custom Location.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.ExtendedLocation`: 2021-08-15
