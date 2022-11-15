---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_service_plan"
description: |-
  Gets information about an existing Service Plan.
---

# Data Source: azurerm_service_plan

Use this data source to access information about an existing Service Plan.

## Example Usage

```hcl
data "azurerm_service_plan" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_service_plan.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Service Plan.

* `resource_group_name` - (Required) The name of the Resource Group where the Service Plan exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Plan.

* `app_service_environment_id` - The ID of the App Service Environment this Service Plan is part of.

* `kind` - A string representing the Kind of Service Plan.

* `location` - The Azure Region where the Service Plan exists.

* `maximum_elastic_worker_count` - The maximum number of workers in use in an Elastic SKU Plan.

* `worker_count` - The number of Workers (instances) allocated.

* `os_type` - The O/S type for the App Services hosted in this plan.

* `per_site_scaling_enabled` - Is Per Site Scaling be enabled?

* `reserved` - Whether this is a reserved Service Plan Type. `true` if `os_type` is `Linux`, otherwise `false`.

* `sku_name` - The SKU for the Service Plan.

* `zone_balancing_enabled` - Is the Service Plan balance across Availability Zones in the region?

* `tags` - A mapping of tags assigned to the Service Plan.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Service Plan.
