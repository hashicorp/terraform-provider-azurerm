---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_plan"
sidebar_current: "docs-azurerm-datasource-app-service-plan"
description: |-
  Get information about an App Service Plan.
---

# Data Source: azurerm_app_service_plan

Use this data source to obtain information about an App Service Plan (formerly known as a `Server Farm`).

## Example Usage

```hcl
data "azurerm_app_service_plan" "test" {
  name                = "search-app-service-plan"
  resource_group_name = "search-service"
}

output "app_service_plan_id" {
  value = "${data.azurerm_app_service_plan.test.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the App Service Plan.
* `resource_group_name` - (Required) The Name of the Resource Group where the App Service Plan exists.

## Attributes Reference

* `id` - The ID of the App Service Plan.

* `location` - The Azure location where the App Service Plan exists

* `kind` - The Operating System type of the App Service Plan

* `sku` - A `sku` block as documented below.

* `properties` - A `properties` block as documented below.

* `tags` - A mapping of tags assigned to the resource.

* `maximum_number_of_workers` - The maximum number of workers supported with the App Service Plan's sku.

---

A `sku` block supports the following:

* `tier` - Specifies the plan's pricing tier.

* `size` - Specifies the plan's instance size.

* `capacity` - Specifies the number of workers associated with this App Service Plan.


A `properties` block supports the following:

* `app_service_environment_id` - The ID of the App Service Environment where the App Service Plan is located.

* `maximum_number_of_workers` - Maximum number of instances that can be assigned to this App Service plan.

* `reserved` - Is this App Service Plan `Reserved`?

* `per_site_scaling` - Can Apps assigned to this App Service Plan be scaled independently?
