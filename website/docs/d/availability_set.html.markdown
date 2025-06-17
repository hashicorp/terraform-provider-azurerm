---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_availability_set"
description: |-
  Gets information about an existing Availability Set.
---

# Data Source: azurerm_availability_set

Use this data source to access information about an existing Availability Set.

## Example Usage

```hcl
data "azurerm_availability_set" "example" {
  name                = "tf-appsecuritygroup"
  resource_group_name = "my-resource-group"
}

output "availability_set_id" {
  value = data.azurerm_availability_set.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Availability Set.

* `resource_group_name` - The name of the resource group in which the Availability Set exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Availability Set.

* `location` - The supported Azure location where the Availability Set exists.

* `managed` - Whether the availability set is managed or not.

* `platform_fault_domain_count` - The number of fault domains that are used.

* `platform_update_domain_count` - The number of update domains that are used.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Availability Set.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Compute`: 2024-03-01
