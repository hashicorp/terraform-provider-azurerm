---
subcategory: "Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_firewall_rule"
description: |-
  Lists Firewall Rules associated with a Redis Cache.
---

# List resource: azurerm_redis_firewall_rule

Lists Firewall Rules associated with a Redis Cache.

## Example Usage

### List all firewall rules deployed in a specific Redis Cache

```hcl
list "azurerm_redis_firewall_rule" "example" {
  provider = azurerm
  config {
    redis_cache_id = "some-redis-cache-id"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `redis_cache_id` - (Required) The full ID of an existing Azure Redis Cache.

````
