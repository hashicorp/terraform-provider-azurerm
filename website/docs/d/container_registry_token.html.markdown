---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_token"
description: |-
  Get information about an existing Container Registry token

---

# Data Source: azurerm_container_registry_token

Use this data source to access information about an existing Container Registry token.

## Example Usage

```hcl
data "azurerm_container_registry_token" "example" {
  name                    = "exampletoken"
  resource_group_name     = "example-resource-group"
  container_registry_name = "example-registry"
}

output "scope_map_id" {
  value = data.azurerm_container_registry_token.example.scope_map_id
}
```

## Argument Reference

* `name` - The name of the Container Registry token.
* `container_registry_name` - The Name of the Container Registry where the token exists.
* `resource_group_name` - The Name of the Resource Group where this Container Registry token exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Container Registry token ID.

* `scope_map_id` - The Scope Map ID used by the token.

* `enabled` - Whether this Token is enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry token.
