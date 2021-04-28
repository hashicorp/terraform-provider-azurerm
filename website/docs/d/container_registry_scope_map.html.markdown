---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_scope_map"
description: |-
  Get information about an existing Container Registry scope map

---

# Data Source: azurerm_container_registry_scope_map

Use this data source to access information about an existing Container Registry scope map.

## Example Usage

```hcl
data "azurerm_container_registry_scope_map" "example" {
  name                    = "example-scope-map"
  resource_group_name     = "example-resource-group"
  container_registry_name = "example-registry"
}

output "actions" {
  value = data.azurerm_container_registry_scope_map.example.actions
}
```

## Argument Reference

* `name` - The name of the Container Registry token.
* `container_registry_name` - The Name of the Container Registry where the token exists.
* `resource_group_name` - The Name of the Resource Group where this Container Registry token exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Container Registry scope map ID.

* `actions` - The actions for the Scope Map.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry scope map.
