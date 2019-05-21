---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall"
sidebar_current: "docs-azurerm-datasource-firewall"
description: |-
  Gets information about an existing Azure Firewall.

---

# Data Source: azurerm_redis_cache

Use this data source to access information about an existing Redis Cache

# Example Usage
```hcl
data "azurerm_redis_cache" "test" {
  name = "myrediscache"
  resource_group_name = "redis-cache"
}

output "primary_access_key" {
  value = "${data.azurerm_redis_cache.test.primary_access_key}"
}

output "hostname" {
  value = "${data.azurerm_redis_cache.test.hostname}"
}

```

## Argument Reference

* `name` - (Required) Specifies the name of the Redis cache
* `resource_group_name` - (Required) Specifies the name of the resource group the Redis cache instance is located in 

## Attribute Reference

* `id` - The Route ID.

* `hostname` - The Hostname of the Redis Instance

* `ssl_port` - The SSL Port of the Redis Instance

* `port` - The non-SSL Port of the Redis Instance

* `primary_access_key` - The Primary Access Key for the Redis Instance

* `secondary_access_key` - The Secondary Access Key for the Redis Instance

* `redis_configuration` - A `redis_configuration` block