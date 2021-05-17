---
subcategory: "Redis Enterprise"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_enterprise_database"
description: |-
  Gets information about an existing Redis Enterprise Database.

---

# Data Source: azurerm_redis_enterprise_database

Use this data source to access information about an existing Redis Enterprise Database

## Example Usage

```hcl
data "azurerm_redis_enterprise_database" "example" {
  name                = "default"
  resource_group_name = azurerm_resource_group.example.name
  cluster_id          = azurerm_redis_enterprise_cluster.example.id
}

output "redis_enterprise_database_primary_key" {
  value       = data.azurerm_redis_enterprise_database.example.primary_access_key
  description = "The Redis Enterprise DB primary key."
}

output "redis_enterprise_database_secondary_key" {
  value       = data.azurerm_redis_enterprise_database.example.secondary_access_key
  description = "The Redis Enterprise DB secondary key."
}
```

## Argument Reference

* `name` - The name of the Redis Enterprise Database.

* `resource_group_name` - The name of the resource group the Redis Enterprise Database instance is located in.

* `cluster_id` - The resource ID of Redis Enterprise Cluster which hosts the Redis Enterprise Database instance.

## Attribute Reference

* `id` - The Redis Enterprise Database ID.

* `name` - The Redis Enterprise Database name.

* `cluster_id` - The Redis Enterprise Cluster ID that is hosting the Redis Enterprise Database.

* `primary_access_key` - The Primary Access Key for the Redis Enterprise Database instance.

* `secondary_access_key` - The Secondary Access Key for the Redis Enterprise Database instance.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Enterprise Database.
