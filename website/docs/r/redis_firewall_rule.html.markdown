---
subcategory: "Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_firewall_rule"
description: |-
  Manages a Firewall Rule associated with a Redis Cache.

---

# azurerm_redis_firewall_rule

Manages a Firewall Rule associated with a Redis Cache.

## Example Usage

```hcl
resource "random_id" "server" {
  keepers = {
    azi_id = 1
  }

  byte_length = 8
}

resource "azurerm_resource_group" "example" {
  name     = "redis-resourcegroup"
  location = "West Europe"
}

resource "azurerm_redis_cache" "example" {
  name                = "redis${random_id.server.hex}"
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

resource "azurerm_redis_firewall_rule" "example" {
  name                = "someIPrange"
  redis_cache_name    = azurerm_redis_cache.example.name
  resource_group_name = azurerm_resource_group.example.name
  start_ip            = "1.2.3.4"
  end_ip              = "2.3.4.5"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Firewall Rule. Changing this forces a new resource to be created.

* `redis_cache_name` - (Required) The name of the Redis Cache. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which this Redis Cache exists. Changing this forces a new resource to be created.

* `start_ip` - (Required) The lowest IP address included in the range

* `end_ip` - (Required) The highest IP address included in the range.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Redis Firewall Rule.

## Timeouts

 The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Redis Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Redis Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Redis Firewall Rule.

## Import

Redis Firewall Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_firewall_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redis/cache1/firewallRules/rule1
```
