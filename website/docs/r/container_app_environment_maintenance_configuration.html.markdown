---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_maintenance_configuration"
description: |-
  Manages a Container App Environment Maintenance Configuration.
---

# azurerm_container_app_environment_maintenance_configuration

Manages a Container App Environment Maintenance Configuration.

~> **Note:** Planned maintenance is a paid feature. For more information, see [Azure Container Apps planned maintenance](https://learn.microsoft.com/en-us/azure/container-apps/planned-maintenance).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-01"
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

resource "azurerm_container_app_environment_maintenance_configuration" "example" {
  container_app_environment_id = azurerm_container_app_environment.example.id

  maintenance_window {
    day_of_week       = "Sunday"
    start_hour_in_utc = 1
    duration_hours    = 8
  }
}
```

## Arguments Reference

The following arguments are supported:

* `container_app_environment_id` - (Required) The ID of the Container App Environment to which this Maintenance Configuration belongs. Changing this forces a new resource to be created.

* `maintenance_window` - (Required) A `maintenance_window` block as defined below.

---

A `maintenance_window` block supports the following:

* `day_of_week` - (Required) The day of the week for the maintenance window. Possible values are `Friday`, `Monday`, `Saturday`, `Sunday`, `Thursday`, `Tuesday`, and `Wednesday`.

* `start_hour_in_utc` - (Required) The start hour of the maintenance window in UTC. Possible values are between `0` and `23`.

* `duration_hours` - (Required) The duration of the maintenance window in hours. Possible values are between `8` and `24`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment Maintenance Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment Maintenance Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Maintenance Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment Maintenance Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment Maintenance Configuration.

## Import

A Container App Environment Maintenance Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment_maintenance_configuration.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myEnvironment/maintenanceConfigurations/default"
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.App` - 2025-07-01
