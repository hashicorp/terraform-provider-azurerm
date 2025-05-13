---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment"
description: |-
  Gets information about a Container App Environment.
---

# Data Source: azurerm_container_app_environment.

Use this data source to access information about an existing Container App Environment.

## Example Usage

```hcl
data "azurerm_container_app_environment" "example" {
  name                = "example-environment"
  resource_group_name = "example-resources"
}
```


## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container Apps Managed Environment.

* `resource_group_name` - (Required) The name of the Resource Group where this Container App Environment exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment.

* `infrastructure_subnet_id` - The ID of the Subnet in use by the Container Apps Control Plane.

~> **Note:** This will only be populated for Environments that have `internal_load_balancer_enabled` set to true.

* `custom_domain_verification_id` - The ID of the Custom Domain Verification for this Container App Environment.

* `default_domain` - The default publicly resolvable name of this Container App Environment. This is generated at creation time to be globally unique.

* `docker_bridge_cidr` - The network addressing in which the Container Apps in this Container App Environment will reside in CIDR notation.

~> **Note:** This will only be populated for Environments that have `internal_load_balancer_enabled` set to true.

* `internal_load_balancer_enabled` - Does the Container App Environment operate in Internal Load Balancing Mode?

* `location` - The Azure Location where this Container App Environment exists.

* `log_analytics_workspace_name` - The name of the Log Analytics Workspace this Container Apps Managed Environment is linked to.

~> **Note:** This will only be populated for Environments that have `logs_destination` set to `log-analytics` and the Log Analytics Workspace is in the same subscription as the Environment.

* `platform_reserved_cidr` - The IP range, in CIDR notation, that is reserved for environment infrastructure IP addresses.

~> **Note:** This will only be populated for Environments that have `internal_load_balancer_enabled` set to true.

* `platform_reserved_dns_ip_address` - The IP address from the IP range defined by `platform_reserved_cidr` that is reserved for the internal DNS server.

~> **Note:** This will only be populated for Environments that have `internal_load_balancer_enabled` set to true.

* `static_ip_address` - The Static IP address of the Environment.

~> **Note:** If `internal_load_balancer_enabled` is true, this will be a Private IP in the subnet, otherwise this will be allocated a Public IPv4 address.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment.
