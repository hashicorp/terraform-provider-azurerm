---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_dapr_component"
description: |-
  Manages a Dapr Component for a Container App Environment.
---

# azurerm_container_app_environment_dapr_component

Manages a Dapr Component for a Container App Environment.

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
  name                       = "Example-Environment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_container_app_environment_dapr_component" "example" {
  name                         = "example-component"
  container_app_environment_id = azurerm_container_app_environment.example.id
  component_type               = "state.azure.blobstorage"
  version                      = "v1"
}
```

## Arguments Reference

The following arguments are supported:

* `container_app_environment_id` - (Required) The ID of the Container App Managed Environment for this Dapr Component. Changing this forces a new resource to be created.

* `name` - (Required) The name for this Dapr component. Changing this forces a new resource to be created.

* `component_type` - (Required) The Dapr Component Type. For example `state.azure.blobstorage`. Changing this forces a new resource to be created.

* `version` - (Required) The version of the component.

---

* `ignore_errors` - (Optional) Should the Dapr sidecar to continue initialisation if the component fails to load. Defaults to `false`

* `init_timeout` - (Optional) The timeout for component initialisation as a `ISO8601` formatted string. e.g. `5s`, `2h`, `1m`. Defaults to `5s`.

* `metadata` - (Optional) One or more `metadata` blocks as detailed below.

* `scopes` - (Optional) A list of scopes to which this component applies.

~> **Note:** See the official docs for more information at https://learn.microsoft.com/en-us/azure/container-apps/dapr-overview?tabs=bicep1%2Cyaml#component-scopes

* `secret` - (Optional) A `secret` block as detailed below.

---

A `metadata` block supports the following:

* `name` - (Required) The name of the Metadata configuration item.

* `secret_name` - (Optional) The name of a secret specified in the `secrets` block that contains the value for this metadata configuration item.

* `value` - (Optional) The value for this metadata configuration item.

---

A `secret` block supports the following:

* `name` - (Required) The Secret name.

* `value` - (Required) The value for this secret.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment Dapr Component


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment Dapr Component.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Dapr Component.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment Dapr Component.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment Dapr Component.

## Import

A Dapr Component for a Container App Environment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment_dapr_component.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myenv/daprComponents/mydaprcomponent"
```
