---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_dapr_component"
description: |-
  Gets information about a Dapr Component in a Container App Environment.
---

# Data Source: azurerm_container_app_environment_dapr_component.

## Example Usage

```hcl
data "azurerm_container_app_environment" "example" {
  name                = "example-environment"
  resource_group_name = "example-resources"
}

resource "azurerm_container_app_environment_dapr_component" "example" {
  name                         = "example-component"
  container_app_environment_id = data.azurerm_container_app_environment.example.id

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name for this Dapr component. Changing this forces a new resource to be created.

* `container_app_environment_id` - (Required) The Container App Managed Environment ID to configure this Dapr component on. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment Dapr Component

* `ignore_errors` - Should the Dapr sidecar to continue initialisation if the component fails to load. Defaults to `false`

* `init_timeout` - The component initialisation timeout in ISO8601 format. e.g. `5s`, `2h`, `1m`. Defaults to `5s`

* `metadata` - A `metadata` block as detailed below.

* `scopes` - A list of scopes to which this component applies. e.g. a Container App's `dapr.app_id` value.

* `secret` - A `secret` block as detailed below.

* `type` - The Dapr Component Type.

* `version` - The version of the component.

---

A `metadata` block exports the following:

* `name` -  The name of the Metadata configuration item.

* `secret_name` -  The name of a secret specified in the `secrets` block that contains the value for this metadata configuration item.

* `value` -  The value for this metadata configuration item.

---

A `secret` block exports the following:

* `name` -  The secret name.

* `value` -  The value for this secret.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Dapr Component.
