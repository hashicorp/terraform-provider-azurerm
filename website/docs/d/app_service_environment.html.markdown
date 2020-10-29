---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_app_service_environment"
description: |-
  Gets information about an existing App Service Environment.
---

# Data Source: azurerm_app_service_environment

Use this data source to access information about an existing App Service Environment.

## Example Usage

```hcl
data "azurerm_app_service_environment" "example" {
  name                = "existing-ase"
  resource_group_name = "existing-rg"
}

output "id" {
  value = data.azurerm_app_service_environment.example.id
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of this App Service Environment.

- `resource_group_name` - (Required) The name of the Resource Group where the App Service Environment exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the App Service Environment.

- `front_end_scale_factor` - The number of app instances per App Service Environment Front End.

- `internal_ip_address` - IP address of internal load balancer of the App Service Environment.

- `location` - The Azure Region where the App Service Environment exists.

- `outbound_ip_addresses` - List of outbound IP addresses of the App Service Environment.

- `pricing_tier` - The Pricing Tier (Isolated SKU) of the App Service Environment.

- `service_ip_address` - IP address of service endpoint of the App Service Environment.

- `tags` - A mapping of tags assigned to the App Service Environment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

- `read` - (Defaults to 5 minutes) Used when retrieving the App Service Environment.
