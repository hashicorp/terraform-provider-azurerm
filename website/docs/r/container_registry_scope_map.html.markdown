---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_scope_map"
description: |-
  Manages an Azure Container Registry scope map.

---

# azurerm_container_registry_scope_map

Manages an Azure Container Registry scope map.  For more information on scope maps see the [product documentation](https://learn.microsoft.com/en-us/azure/container-registry/container-registry-repository-scoped-permissions).

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_container_registry" "example" {
  name                = "exampleregistry"
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the scope map. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container Registry token. Changing this forces a new resource to be created.

* `container_registry_name` - (Required) The name of the Container Registry. Changing this forces a new resource to be created.

* `actions` - (Required) A list of actions to attach to the scope map (e.g. `repo/content/read`, `repo2/content/delete`).

* `description` - (Optional) The description of the Container Registry.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Registry scope map.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry scope map.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry scope map.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry scope map.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry scope map.

## Import

Container Registries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_scope_map.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1/scopeMaps/scopemap1
```
