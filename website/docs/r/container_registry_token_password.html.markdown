---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_token_password"
description: |-
  Manages a Container Registry Token Password.
---

# azurerm_container_registry_token_password

Manages a Container Registry Token Password associated with a scope map.  For more information on scope maps and their tokens see the [product documentation](https://learn.microsoft.com/en-us/azure/container-registry/container-registry-repository-scoped-permissions).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_container_registry" "example" {
  name                     = "example-registry"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  sku                      = "Basic"
  admin_enabled            = false
  georeplication_locations = ["East US", "West Europe"]
}

resource "azurerm_container_registry_scope_map" "example" {
  name                    = "example-scope-map"
  container_registry_name = azurerm_container_registry.example.name
  resource_group_name     = azurerm_resource_group.example.name
  actions = [
    "repositories/repo1/content/read",
    "repositories/repo1/content/write"
  ]
}

resource "azurerm_container_registry_token" "example" {
  name                    = "exampletoken"
  container_registry_name = azurerm_container_registry.example.name
  resource_group_name     = azurerm_resource_group.example.name
  scope_map_id            = azurerm_container_registry_scope_map.example.id
}

resource "azurerm_container_registry_token_password" "example" {
  container_registry_token_id = azurerm_container_registry_token.example.id

  password1 {
    expiry = "2023-03-22T17:57:36+08:00"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `container_registry_token_id` - (Required) The ID of the Container Registry Token that this Container Registry Token Password resides in. Changing this forces a new Container Registry Token Password to be created.

* `password1` - (Required) One `password` block as defined below.

* `password2` - (Optional) One `password` block as defined below.

---

A `password` block supports the following:

* `expiry` - (Optional) The expiration date of the password in RFC3339 format. If not specified, the password never expires. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Registry Token Password.

* `password1` - One `password` block as defined below.

* `password2` - One `password` block as defined below.

---

A `password` block exports the following:

* `value` - The value of the password (Sensitive).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry Token Password.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry Token Password.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry Token Password.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry Token Password.

## Import

Container Registry Token Passwords can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_token_password.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/tokens/token1/passwords/password
```
