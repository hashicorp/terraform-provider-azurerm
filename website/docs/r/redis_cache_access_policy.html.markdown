---
subcategory: "Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_cache_access_policy"
description: |-
  Manages a Redis Cache Access Policy.
---

# azurerm_redis_cache_access_policy

Manages a Redis Cache Access Policy

## Example Usage

```hcl
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

resource "azurerm_redis_cache_access_policy" "example" {
  name           = "example"
  redis_cache_id = azurerm_redis_cache.example.id
  permissions    = "+@read +@connection +cluster|info"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Redis Cache Access Policy. Changing this forces a new Redis Cache Access Policy to be created.

* `redis_cache_id` - (Required) The ID of the Redis Cache. Changing this forces a new Redis Cache Access Policy to be created.

* `permissions` - (Required) Permissions that are going to be assigned to this Redis Cache Access Policy.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Redis Cache Access Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Redis Cache Access Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Cache Access Policy.
* `update` - (Defaults to 5 minutes) Used when updating the Redis Cache Access Policy.
* `delete` - (Defaults to 5 minutes) Used when deleting the Redis Cache Access Policy.

## Import

Redis Cache Access Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_cache_access_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redis/cache1/accessPolicies/policy1
```
