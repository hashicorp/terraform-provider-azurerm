---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_firewall_rule"
sidebar_current: "docs-azurerm-resource-redis-firewall-rule"
description: |-
  Manages a Firewall Rule associated with a Premium Redis Cache.

---

# azurerm_redis_firewall_rule

Manages a Firewall Rule associated with a Premium Redis Cache.

~> **Note:** Redis Firewall Rules can only be assigned to a Redis Cache with a `Premium` SKU.

## Example Usage

```hcl
resource "random_id" "server" {
  keepers = {
    azi_id = 1
  }

  byte_length = 8
}

resource "azurerm_resource_group" "test" {
  name     = "redis-resourcegroup"
  location = "West Europe"
}

resource "azurerm_redis_cache" "test" {
  name                = "redis${random_id.server.hex}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  capacity            = 1
  family              = "P"
  sku_name            = "Premium"
  enable_non_ssl_port = false

  redis_configuration {
    maxclients         = 256
    maxmemory_reserved = 2
    maxmemory_delta    = 2
    maxmemory_policy   = "allkeys-lru"
  }
}

resource "azurerm_redis_firewall_rule" "test" {
  name                = "someIPrange"
  redis_cache_name    = "${azurerm_redis_cache.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  start_ip            = "1.2.3.4"
  end_ip              = "2.3.4.5"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Firewall Rule. Changing this forces a new resource to be created.

* `redis_cache_name` - (Required) The name of the Redis Cache. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which this Redis Cache exists.

* `start_ip` - (Required) The lowest IP address included in the range

* `end_ip` - (Required) The highest IP address included in the range.


## Attributes Reference

The following attributes are exported:

* `id` - The Redis Firewall Rule ID.
