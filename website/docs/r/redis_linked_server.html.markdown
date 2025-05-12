---
subcategory: "Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_linked_server"
description: |-
  Manages a Redis Linked Server.
---

# azurerm_redis_linked_server

Manages a Redis Linked Server (ie Geo Location)

## Example Usage

```hcl
resource "azurerm_resource_group" "example-primary" {
  name     = "example-resources-primary"
  location = "East US"
}

resource "azurerm_redis_cache" "example-primary" {
  name                = "example-cache1"
  location            = azurerm_resource_group.example-primary.location
  resource_group_name = azurerm_resource_group.example-primary.name
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

resource "azurerm_resource_group" "example-secondary" {
  name     = "example-resources-secondary"
  location = "West US"
}

resource "azurerm_redis_cache" "example-secondary" {
  name                = "example-cache2"
  location            = azurerm_resource_group.example-secondary.location
  resource_group_name = azurerm_resource_group.example-secondary.name
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

resource "azurerm_redis_linked_server" "example-link" {
  target_redis_cache_name     = azurerm_redis_cache.example-primary.name
  resource_group_name         = azurerm_redis_cache.example-primary.resource_group_name
  linked_redis_cache_id       = azurerm_redis_cache.example-secondary.id
  linked_redis_cache_location = azurerm_redis_cache.example-secondary.location
  server_role                 = "Secondary"
}
```

## Arguments Reference

The following arguments are supported:

* `linked_redis_cache_id` - (Required) The ID of the linked Redis cache. Changing this forces a new Redis to be created.

* `linked_redis_cache_location` - (Required) The location of the linked Redis cache. Changing this forces a new Redis to be created.

* `target_redis_cache_name` - (Required) The name of Redis cache to link with. Changing this forces a new Redis to be created. (eg The primary role)

* `resource_group_name` - (Required) The name of the Resource Group where the Redis caches exists. Changing this forces a new Redis to be created.

* `server_role` - (Required) The role of the linked Redis cache (eg "Secondary"). Changing this forces a new Redis to be created. Possible values are `Primary` and `Secondary`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Redis.

* `name` - The name of the linked server.

* `geo_replicated_primary_host_name` - The geo-replicated primary hostname for this linked server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Redis.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis.
* `delete` - (Defaults to 1 hour) Used when deleting the Redis.

## Import

Redis can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_linked_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redis/cache1/linkedServers/cache2
```
