---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_plan"
description: |-
  Gets information about an existing App Service Plan.
---

# Data Source: azurerm_app_service_plan

Use this data source to access information about an existing App Service Plan (formerly known as a `Server Farm`).

!> **Note:** The `azurerm_app_service_plan` data source is deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use the [`azurerm_service_plan`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/service_plan) data source instead.

## Example Usage

```hcl
data "azurerm_app_service_plan" "example" {
  name                = "search-app-service-plan"
  resource_group_name = "search-service"
}

output "app_service_plan_id" {
  value = data.azurerm_app_service_plan.example.id
}
```

## Argument Reference

* `name` - The name of the App Service Plan.
* `resource_group_name` - The Name of the Resource Group where the App Service Plan exists.

## Attributes Reference

* `id` - The ID of the App Service Plan.

* `location` - The Azure location where the App Service Plan exists

* `kind` - The Operating System type of the App Service Plan

* `sku` - A `sku` block as documented below.

* `app_service_environment_id` - The ID of the App Service Environment where the App Service Plan is located.

* `maximum_number_of_workers` - Maximum number of instances that can be assigned to this App Service plan.

* `reserved` - Is this App Service Plan `Reserved`?

* `per_site_scaling` - Can Apps assigned to this App Service Plan be scaled independently?

* `tags` - A mapping of tags assigned to the resource.

* `maximum_elastic_worker_count` - The maximum number of total workers allowed for this ElasticScaleEnabled App Service Plan.

* `is_xenon` - A flag that indicates if it's a xenon plan (support for Windows Container)

* `maximum_number_of_workers` - The maximum number of workers supported with the App Service Plan's sku.

* `zone_redundant` - App Service Plan perform availability zone balancing.

---

A `sku` block supports the following:

* `tier` - Specifies the plan's pricing tier.

* `size` - Specifies the plan's instance size.

* `capacity` - Specifies the number of workers associated with this App Service Plan.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Plan.
