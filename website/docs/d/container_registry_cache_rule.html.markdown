---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_cache_rule"
description: |-
  Get information about an existing Container Registry Cache Rule

---

# Data Source: azurerm_container_registry_cache_rule

Use this data source to access information about an existing Container Registry Cache Rule.

## Example Usage

```hcl
data "azurerm_container_registry" "example" {
  name                  = "testacr"
  container_registry_id = "test"
}

output "cache_rule_source_repo" {
  value = data.azurerm_container_registry_cache_rule.example.source_repo
}
```

## Argument Reference

* `name` - Specifies the name of the Container Registry Cache Rule. Only Alphanumeric characters allowed. Changing this forces a new resource to be created.

* `container_registry_id` - The ID of the container registry where the cache rule should apply. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Container Registry ID.

* `source_repo` - The name of the source repository path.

* `target_repo` - The name of the new repository path to store artifacts.

* `credential_set_id` - The ARM resource ID of the credential store which is associated with the cache rule.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry.
