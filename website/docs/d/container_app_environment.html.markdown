---
subcategory: "Containerapps"
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
  name                = "exampleContainerAppEnvironment"
  resource_group_name = "exampleResourceGroup"
}
```


## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container Apps Managed Environment.

* `resource_group_name` - (Required) The name of the resource group in which the Container App Environment is to be found.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment

* `apps_subnet_id` - The existing Subnet in use by the Container Apps runtime.

* `control_plane_subnet_id` - The existing Subnet in use by the Container Apps Control Plane.

* `created_at` - The time and date at which this Container App Environment was created.

* `created_by` - The user or principal which created this Container App Environment.

* `created_by_type` - The type of account which created this Container App Environment.

* `default_domain` - The default publicly resolvable name of this Container App Environment

* `docker_bridge_cidr` - The network addressing in which the Container Apps in this Container App Environment will reside in CIDR notation.

* `internal_load_balancer_enabled` - Does the Container Environment operate in Internal Load Balancing Mode?

* `last_modified_at` - The time and date at which this Container App Environment was last modified.

* `last_modified_by` - The user or principal which last modified this Container App Environment.

* `last_modified_by_type` - The type of account which last modified this Container App Environment.

* `location` - The Azure Location in which this resource resides.

* `log_analytics_workspace_name` - The name of the Log Analytics Workspace this Container Apps Managed Environment is linked to.

* `platform_reserved_cidr` - The IP range, in CIDR notation, that is reserved for environment infrastructure IP addresses.

* `platform_reserved_dns_ip` - The IP address from the IP range defined by `platform_reserved_cidr` that is reserved for the internal DNS server.

* `static_ip` - The Static IP of the Environment.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment.
