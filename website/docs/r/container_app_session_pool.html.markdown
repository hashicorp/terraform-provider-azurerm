---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_session_pool"
description: |-
  Manages a Container App Session Pool.
---

# azurerm_container_app_session_pool

Manages a Container App Session Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_app_session_pool" "example" {
  name                       = "examplesessionpool"
  resource_group_name        = azurerm_resource_group.example.name
  location                   = azurerm_resource_group.example.location
  cooldown_period_in_seconds = 300
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Container App Session Pool. Changing this forces a new resource to be created.

~> **Note:** `name` must be between 3 and 63 characters long, begin with a lowercase letter and contain only lowercase letters and numbers.

* `resource_group_name` - (Required) The name of the Resource Group where the Container App Session Pool should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Container App Session Pool should exist. Changing this forces a new resource to be created.

* `container_type` - (Optional) The type of container used for the sessions in this Container App Session Pool. Possible values are `CustomContainer` and `PythonLTS`. Defaults to `PythonLTS`. Changing this forces a new resource to be created.

* `max_concurrent_sessions` - (Optional) The maximum number of sessions which can run concurrently in this Container App Session Pool. Defaults to `5`.

* `container_app_environment_id` - (Optional) The ID of the Container App Environment used to host the sessions in this Container App Session Pool. Changing this forces a new resource to be created.

~> **Note:** `container_app_environment_id` is required when `container_type` is set to `CustomContainer`.

~> **Note:** The Container App Environment must be of type `Workload profile`. `Consumption only` environments are not supported.

* `cooldown_period_in_seconds` - (Optional) The number of seconds a session remains alive once it becomes idle. Possible values range between `300` and `3600`.

~> **Note:** `cooldown_period_in_seconds` must be specified when `lifecycle_type` is set to `Timed`, cannot be specified otherwise, and conflicts with `max_alive_period_in_seconds`.

* `custom_container_template` - (Optional) A `custom_container_template` block as defined below.

~> **Note:** `custom_container_template` must be specified when `container_type` is set to `CustomContainer`, and may not be specified otherwise.

* `identity` - (Optional) An `identity` block as defined below.

~> **Note:** `identity` may only be specified when `container_type` is set to `CustomContainer`.

* `lifecycle_type` - (Optional) The lifecycle type of the sessions in this Container App Session Pool. Possible values are `OnContainerExit` and `Timed`. Defaults to `Timed`.

~> **Note:** `OnContainerExit` can only be used when `container_type` is set to `CustomContainer`.

* `max_alive_period_in_seconds` - (Optional) The maximum number of seconds a session remains alive.

~> **Note:** `max_alive_period_in_seconds` must be specified when `lifecycle_type` is set to `OnContainerExit`, cannot be specified otherwise, and conflicts with `cooldown_period_in_seconds`.

* `network_egress_enabled` - (Optional) Should sessions in this Container App Session Pool be able to make outbound network requests? Defaults to `false`.

* `ready_session_instances` - (Optional) The minimum number of sessions which are kept ready in this Container App Session Pool. This must be greater than `0` and less than `max_concurrent_sessions`.

~> **Note:** `ready_session_instances` must be specified when `container_type` is set to `CustomContainer`.

* `secret` - (Optional) One or more `secret` blocks as defined below.

* `session_managed_identities` - (Optional) A list of Managed Identities which are made available to the code running inside the sessions. Each value should be the ID of a User Assigned Managed Identity, or `System` to use the System Assigned Identity.

~> **Note:** Each Managed Identity listed here must also be assigned to this Container App Session Pool via the `identity` block. Listing it here grants the code running in a session access to that Managed Identity.

* `tags` - (Optional) A mapping of tags which should be assigned to the Container App Session Pool.

---

A `container` block supports the following:

* `image` - (Required) The image used to create the container.

* `name` - (Required) The name which should be used for this container.

* `args` - (Optional) A list of arguments to pass to the container.

* `command` - (Optional) A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.

* `cpu` - (Optional) The amount of vCPU to allocate to the container, e.g. `0.25`.

* `env` - (Optional) One or more `env` blocks as defined below.

* `memory` - (Optional) The amount of memory to allocate to the container, e.g. `0.5Gi`.

---

A `custom_container_template` block supports the following:

* `container` - (Required) One or more `container` blocks as defined below.

* `ingress_target_port` - (Required) The target port in the container which receives traffic from the ingress. Possible values range between `1` and `65535`.

* `registry` - (Optional) A `registry` block as defined below.

---

An `env` block supports the following:

* `name` - (Required) The name of the environment variable.

* `secret_name` - (Optional) The name of the `secret` which contains the value for this environment variable.

* `value` - (Optional) The value for this environment variable.

~> **Note:** `value` is ignored when `secret_name` is specified.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Identity assigned to this Container App Session Pool. Possible values are `SystemAssigned`, `UserAssigned`, and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Container App Session Pool.

~> **Note:** `identity_ids` is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `registry` block supports the following:

* `server` - (Required) The hostname of the Container Registry.

* `identity` - (Optional) The ID of a User Assigned Managed Identity, or `System` to use the System Assigned Identity, used to pull images from the Container Registry.

* `password_secret_name` - (Optional) The name of the `secret` which contains the login password for the Container Registry.

* `username` - (Optional) The username to use for the Container Registry.

~> **Note:** `identity` conflicts with `username` and `password_secret_name`.

---

A `secret` block supports the following:

* `name` - (Required) The name of the secret.

* `value` - (Required) The value for this secret.

~> **Note:** `value` is not returned by the Azure API, so it is stored from the Terraform configuration. Changes made to a secret value outside of Terraform cannot be detected.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Session Pool.

* `identity` - An `identity` block as defined below.

* `node_count` - The number of nodes in use by this Container App Session Pool.

* `pool_management_endpoint` - The endpoint used to manage the sessions in this Container App Session Pool.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID of the System Assigned Managed Identity assigned to this Container App Session Pool.

* `tenant_id` - The Tenant ID of the System Assigned Managed Identity assigned to this Container App Session Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Session Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Session Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Session Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Session Pool.

## Import

A Container App Session Pool can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_session_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.App/sessionPools/sessionPool1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.App` - 2025-07-01
