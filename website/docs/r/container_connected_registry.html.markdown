---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_connected_registry"
description: |-
  Manages a Container Connected Registry.
---

# azurerm_container_connected_registry

Manages a Container Connected Registry.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}
resource "azurerm_container_registry" "example" {
  name                  = "exampleacr"
  resource_group_name   = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  sku                   = "Premium"
  data_endpoint_enabled = true
}
resource "azurerm_container_registry_scope_map" "example" {
  name                    = "examplescopemap"
  container_registry_name = azurerm_container_registry.example.name
  resource_group_name     = azurerm_container_registry.example.resource_group_name
  actions = [
    "repositories/hello-world/content/delete",
    "repositories/hello-world/content/read",
    "repositories/hello-world/content/write",
    "repositories/hello-world/metadata/read",
    "repositories/hello-world/metadata/write",
    "gateway/examplecr/config/read",
    "gateway/examplecr/config/write",
    "gateway/examplecr/message/read",
    "gateway/examplecr/message/write",
  ]
}
resource "azurerm_container_registry_token" "example" {
  name                    = "exampletoken"
  container_registry_name = azurerm_container_registry.example.name
  resource_group_name     = azurerm_container_registry.example.resource_group_name
  scope_map_id            = azurerm_container_registry_scope_map.example.id
}
resource "azurerm_container_connected_registry" "example" {
  name                  = "examplecr"
  container_registry_id = azurerm_container_registry.example.id
  sync_token_id         = azurerm_container_registry_token.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `container_registry_id` - (Required) The ID of the Container Registry that this Connected Registry will reside in. Changing this forces a new Container Connected Registry to be created.

-> **Note:** If `parent_registry_id` is not specified, the Connected Registry will be connected to the Container Registry identified by `container_registry_id`.

* `name` - (Required) The name which should be used for this Container Connected Registry. Changing this forces a new Container Connected Registry to be created.

* `sync_token_id` - (Required) The ID of the Container Registry Token which is used for synchronizing the Connected Registry. Changing this forces a new Container Connected Registry to be created.

---

* `audit_log_enabled` - (Optional) Should the log auditing be enabled?

* `client_token_ids` - (Optional) Specifies a list of IDs of Container Registry Tokens, which are meant to be used by the clients to connect to the Connected Registry.

* `log_level` - (Optional) The verbosity of the logs. Possible values are `None`, `Debug`, `Information`, `Warning` and `Error`. Defaults to `None`.

* `mode` - (Optional) The mode of the Connected Registry. Possible values are `Mirror`, `ReadOnly`, `ReadWrite` and `Registry`. Changing this forces a new Container Connected Registry to be created. Defaults to `ReadWrite`.

* `notification` - (Optional) One or more `notification` blocks as defined below.

* `parent_registry_id` - (Optional) The ID of the parent registry. This can be either a Container Registry ID or a Connected Registry ID. Changing this forces a new Container Connected Registry to be created.

* `sync_message_ttl` - (Optional) The period of time (in form of ISO8601) for which a message is available to sync before it is expired. Allowed range is from `P1D` to `P90D`. Defaults to `P1D`.

* `sync_schedule` - (Optional) The cron expression indicating the schedule that the Connected Registry will sync with its parent. Defaults to `* * * * *`.

* `sync_window` - (Optional) The time window (in form of ISO8601) during which sync is enabled for each schedule occurrence. Allowed range is from `PT3H` to `P7D`.

---

A `notification` block supports the following:

* `name` - (Required) The name of the artifact that wants to be subscribed for the Connected Registry.

* `action` - (Required) The action of the artifact that wants to be subscribed for the Connected Registry. Possible values are `push`, `delete` and `*` (i.e. any).

* `tag` - (Optional) The tag of the artifact that wants to be subscribed for the Connected Registry.

* `digest` - (Optional) The digest of the artifact that wants to be subscribed for the Connected Registry.

~> **Note:** One of either `tag` or `digest` can be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Connected Registry.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Connected Registry.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Connected Registry.
* `update` - (Defaults to 30 minutes) Used when updating the Container Connected Registry.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Connected Registry.

## Import

Container Connected Registries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_connected_registry.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/connectedRegistries/registry1
```
