---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_token"
description: |-
  Manages an Azure Container Registry token.

---

# azurerm_container_registry_token

Manages an Azure Container Registry token associated to a scope map. For more information on scope maps and their tokens see the [product documentation](https://learn.microsoft.com/en-us/azure/container-registry/container-registry-repository-scoped-permissions).

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_container_registry" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Basic"
  admin_enabled       = false

  georeplications {
    location = "East US"
  }
  georeplications {
    location = "West Europe"
  }
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the token. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container Registry token. Changing this forces a new resource to be created.

* `container_registry_name` - (Required) The name of the Container Registry. Changing this forces a new resource to be created.

* `scope_map_id` - (Required) The ID of the Container Registry Scope Map associated with the token.

* `enabled` - (Optional) Should the Container Registry token be enabled? Defaults to `true`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Registry token.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry token.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry token.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry token.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry token.

## Import

Container Registries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_token.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1/tokens/token1
```
