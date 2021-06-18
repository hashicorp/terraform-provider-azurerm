---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_scope_map"
description: |-
  Manages an Azure Container Registry scope map.

---

# azurerm_container_registry_scope_map

Manages an Azure Container Registry scope map.  Scope Maps are a preview feature only available in Premium SKU Container registries.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_container_registry" "example" {
  name                     = "example-registry"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  sku                      = "Premium"
  admin_enabled            = false
  georeplication_locations = ["East US", "West Europe"]
}

resource "azurerm_container_registry_scope_map" "example" {
  name                    = "example-scope-map"
  container_registry_name = azurerm_container_registry.acr.name
  resource_group_name     = azurerm_resource_group.rg.name
  actions = [
    "repositories/repo1/content/read",
    "repositories/repo1/content/create"
  ]
}
```

## Argument Reference

The following arguments are supported:


* `name` - (Required) Specifies the name of the scope map. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container Registry token. Changing this forces a new resource to be created.

* `container_registry_name` - (Required) The name of the Container Registry. Changing this forces a new resource to be created.

* `actions` - (Required) A list of actions to attach to the scope map (e.g. `repo/content/read`, `repo2/content/delete`).

---
## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Container Registry scope map.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry scope map.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry scope map.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry scope map.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry scope map.

## Import

Container Registries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_scope_map.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1/scopeMaps/scopemap1
```
