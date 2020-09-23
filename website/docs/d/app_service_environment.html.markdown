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
data "azurerm_app_service_environment" "example" {
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

* `internal_ip_address` - IP address of internal load balancer of the App Service Environment.

* `service_ip_address` - IP address of service endpoint of the App Service Environment.

* `outbound_ip_addresses` - Outbound IP addresses of the App Service Environment.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Environment.
