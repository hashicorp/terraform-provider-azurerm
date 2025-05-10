---
subcategory: "Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_cache_access_policy_assignment"
description: |-
  Manages a Redis Cache Access Policy Assignment.
---

# azurerm_redis_cache_access_policy_assignment

Manages a Redis Cache Access Policy Assignment

## Example Usage

```hcl
data "azurerm_client_config" "test" {
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_redis_cache" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    maxmemory_reserved = 2
    maxmemory_delta    = 2
    maxmemory_policy   = "allkeys-lru"
  }
}

resource "azurerm_redis_cache_access_policy_assignment" "example" {
  name               = "example"
  redis_cache_id     = azurerm_redis_cache.example.id
  access_policy_name = "Data Contributor"
  object_id          = data.azurerm_client_config.test.object_id
  object_id_alias    = "ServicePrincipal"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Redis Cache Access Policy Assignment. Changing this forces a new Redis Cache Access Policy Assignment to be created.

* `redis_cache_id` - (Required) The ID of the Redis Cache. Changing this forces a new Redis Cache Access Policy Assignment to be created.

* `access_policy_name` - (Required) The name of the Access Policy to be assigned. Changing this forces a new Redis Cache Access Policy Assignment to be created.

* `object_id` - (Required) The principal ID to be assigned the Access Policy. Changing this forces a new Redis Cache Access Policy Assignment to be created.

* `object_id_alias` - (Required) The alias of the principal ID. User-friendly name for object ID. Also represents username for token based authentication. Changing this forces a new Redis Cache Access Policy Assignment to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Redis Cache Access Policy Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Redis Cache Access Policy Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Cache Access Policy Assignment.
* `delete` - (Defaults to 5 minutes) Used when deleting the Redis Cache Access Policy Assignment.

## Import

Redis Cache Policy Assignment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_cache_access_policy_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redis/cache1/accessPolicyAssignments/assignment1
```
