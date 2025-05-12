---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment"
description: |-
  Manages a Container App Environment.
---

# azurerm_container_app_environment

Manages a Container App Environment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "example" {
  name                       = "my-environment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  logs_destination           = "log-analytics"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container Apps Managed Environment. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Container App Environment is to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Container App Environment is to exist. Changing this forces a new resource to be created.

---

* `dapr_application_insights_connection_string` - (Optional) Application Insights connection string used by Dapr to export Service to Service communication telemetry. Changing this forces a new resource to be created.

* `infrastructure_resource_group_name` - (Optional) Name of the platform-managed resource group created for the Managed Environment to host infrastructure resources. Changing this forces a new resource to be created.

~> **Note:** Only valid if a `workload_profile` is specified. If `infrastructure_subnet_id` is specified, this resource group will be created in the same subscription as `infrastructure_subnet_id`.

* `infrastructure_subnet_id` - (Optional) The existing Subnet to use for the Container Apps Control Plane. Changing this forces a new resource to be created. 

~> **Note:** The Subnet must have a `/21` or larger address space.

* `internal_load_balancer_enabled` - (Optional) Should the Container Environment operate in Internal Load Balancing Mode? Defaults to `false`. Changing this forces a new resource to be created.

~> **Note:** can only be set to `true` if `infrastructure_subnet_id` is specified.

* `zone_redundancy_enabled` - (Optional) Should the Container App Environment be created with Zone Redundancy enabled? Defaults to `false`. Changing this forces a new resource to be created.

~> **Note:** can only be set to `true` if `infrastructure_subnet_id` is specified.

* `log_analytics_workspace_id` - (Optional) The ID for the Log Analytics Workspace to link this Container Apps Managed Environment to. 

~> **Note:** required if `logs_destination` is set to `log-analytics`. Cannot be set if `logs_destination` is set to `azure-monitor`.

* `logs_destination` - (Optional) Where the application logs will be saved for this Container Apps Managed Environment. Possible values include `log-analytics` and `azure-monitor`. Omitting this value will result in logs being streamed only.

* `workload_profile` - (Optional) One or more `workload_profile` blocks as defined below.

* `mutual_tls_enabled` - (Optional) Should mutual transport layer security (mTLS) be enabled? Defaults to `false`.

~> **Note:** This feature is in public preview. Enabling mTLS for your applications may increase response latency and reduce maximum throughput in high-load scenarios.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `workload_profile` block supports the following:

* `name` - (Required) The name of the workload profile.

* `workload_profile_type` - (Required) Workload profile type for the workloads to run on. Possible values include `Consumption`, `D4`, `D8`, `D16`, `D32`, `E4`, `E8`, `E16` and `E32`.

~> **Note:** A `Consumption` type must have a name of `Consumption` and an environment may only have one `Consumption` Workload Profile.

~> **Note:** Defining a `Consumption` profile is optional, however, Environments created without an initial Workload Profile cannot have them added at a later time and must be recreated. Similarly, an environment created with Profiles must always have at least one defined Profile, removing all profiles will force a recreation of the resource.

* `maximum_count` - (Required) The maximum number of instances of workload profile that can be deployed in the Container App Environment.

* `minimum_count` - (Required) The minimum number of instances of workload profile that can be deployed in the Container App Environment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment

* `custom_domain_verification_id` - The ID of the Custom Domain Verification for this Container App Environment.

* `default_domain` - The default, publicly resolvable, name of this Container App Environment.

~> **Note:** This value is generated by the service to be globally unique.

* `docker_bridge_cidr` - The network addressing in which the Container Apps in this Container App Environment will reside in CIDR notation.

~> **Note:** This property only has a value when `infrastructure_subnet_id` is configured and will be a range within the CIDR of the Subnet.

* `platform_reserved_cidr` - The IP range, in CIDR notation, that is reserved for environment infrastructure IP addresses.

~> **Note:** This property only has a value when `infrastructure_subnet_id` is configured and will be a range within the CIDR of the Subnet.

* `platform_reserved_dns_ip_address` - The IP address from the IP range defined by `platform_reserved_cidr` that is reserved for the internal DNS server.

~> **Note:** This property only has a value when `infrastructure_subnet_id` is configured and will be a value within the CIDR of the Subnet.

* `static_ip_address` - The Static IP address of the Environment.

~> **Note:** This will be a Public IP unless `internal_load_balancer_enabled` is set to `true`, in which case an IP in the Internal Subnet will be reserved.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment.

## Import

A Container App Environment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myEnvironment"
```
