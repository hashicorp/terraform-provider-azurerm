---
subcategory: "Containerapps"
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
  name                = "acctest-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "example" {
  name                       = "myEnvironment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container Apps Managed Environment. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Container App Environment is to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource is to exist. Changing this forces a new resource to be created.

* `log_analytics_workspace_id` - (Required) The ID for the Log Analytics Workspace to link this Container Apps Managed Environment to. Changing this forces a new resource to be created.

---

* `infrastructure_subnet_id` - (Optional) The existing Subnet to use for the Container Apps Control Plane. **NOTE:** The Subnet must have a `/21` or larger address space. Changing this forces a new resource to be created.

* `internal_load_balancer_enabled` - (Optional) Should the Container Environment operate in Internal Load Balancing Mode? Defaults to `false`. **Note:** can only be set to `true` if `infrastructure_subnet_id` is specified. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment

* `created_at` - The time and date at which this Container App Environment was created.

* `created_by` - The user or principal which created this Container App Environment.

* `created_by_type` - The type of account which created this Container App Environment.

* `default_domain` - The default, publicly resolvable, name of this Container App Environment.

* `docker_bridge_cidr` - The network addressing in which the Container Apps in this Container App Environment will reside in CIDR notation.

* `last_modified_at` - The time and date at which this Container App Environment was last modified.

* `last_modified_by` - The user or principal which last modified this Container App Environment.

* `last_modified_by_type` - The type of account which last modified this Container App Environment.

* `platform_reserved_cidr` - The IP range, in CIDR notation, that is reserved for environment infrastructure IP addresses.

* `platform_reserved_dns_ip` - The IP address from the IP range defined by `platform_reserved_cidr` that is reserved for the internal DNS server.

* `static_ip` - The Static IP of the Environment.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment.

## Import

a Container App Environment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myEnvironment"
```