---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_environment"
description: |-
  Gets information about an existing App Service Environment.
---

# Data Source: azurerm_app_service_environment

Use this data source to access information about an existing App Service Environment

## Example Usage

```hcl
data "azure_app_service_plan" "example" {
  name                = "example-ase"
  resource_group_name = "example-rg"
}

output "app_service_environment_id" {
  value = "${data.azurerm_app_service_environment.id}"
}

```

## Argument Reference

* `name` - (Required) The name of the App Service Environment.

* `resource_group_name` - (Required) The Name of the Resource Group where the App Service Environment exists.

## Attribute Reference

* `id` - The ID of the App Service Environment.

* `location` - The Azure location where the App Service Environment exists

* `front_end_scale_factor` - The number of app instances per App Service Environment Front End

* `pricing_tier` - The Pricing Tier (Isolated SKU) of the App Service Environment.

* `tags` - A mapping of tags assigned to the resource.