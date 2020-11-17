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
  name                = "testtoken"
  resource_group_name = "test"
  container_registry_name = "testacr"
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

* `status` - The status of the token (`enabled` or `disabled`).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry token.
