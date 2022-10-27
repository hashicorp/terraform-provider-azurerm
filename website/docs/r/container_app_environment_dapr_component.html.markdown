---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_dapr_component"
description: |-
  Manages a Container App Environment Dapr Component.
---

# azurerm_container_app_environment_dapr_component

Manages a Container App Environment Dapr Component.

## Example Usage

```hcl
resource "azurerm_container_app_environment_dapr_component" "example" {
  container_app_environment_id = "example"
  name                         = "example"
  type                         = "example"
  version                      = "example"

}
```

## Arguments Reference

The following arguments are supported:

* `container_app_environment_id` - (Required) The Container App Managed Environment ID to configure this Dapr component on. Changing this forces a new resource to be created.

* `name` - (Required) The name for this Dapr component. Changing this forces a new resource to be created.

* `type` - (Required) The Dapr Component Type. For example `state.azure.blobstorage`. Changing this forces a new resource to be created.

* `version` - (Required) The version of the component.

---

* `ignore_errors` - (Optional) Should the Dapr sidecar to continue initialisation if the component fails to load. Defaults to `false`

* `init_timeout` - (Optional) The component initialisation timeout in ISO8601 format. e.g. `5s`, `2h`, `1m`. Defaults to `5s`

* `metadata` - (Optional) A `metadata` block as detailed below.

* `scopes` - (Optional) A list of scopes to which this component applies. e.g. a Container App's `dapr.app_id` value.

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
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment Dapr Component.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Dapr Component.
* `delete` - (Defaults to 5 minutes) Used when deleting the Container App Environment Dapr Component.

## Import

a Container App Environment Dapr Component can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment_dapr_component.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/daprComponents/mydaprcomponent"
```
