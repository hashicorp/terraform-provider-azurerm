---
subcategory: "Managed Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_redis_database"
description: |-
  Gets information about an existing Managed Redis Database.

---

# Data Source: azurerm_managed_redis_database

Use this data source to access information about an existing Managed Redis Database

## Example Usage

```hcl
data "azurerm_managed_redis_database" "example" {
  name       = "default"
  cluster_id = azurerm_managed_redis_cluster.example.id
}

output "managed_redis_database_primary_key" {
  value       = data.azurerm_managed_redis_database.example.primary_access_key
  description = "The Managed Redis DB primary key."
  sensitive   = true
}

output "managed_redis_database_secondary_key" {
  value       = data.azurerm_managed_redis_database.example.secondary_access_key
  description = "The Managed Redis DB secondary key."
  sensitive   = true
}
```

## Argument Reference

* `name` - The name of the Managed Redis Database.

* `cluster_id` - The resource ID of Managed Redis Cluster which hosts the Managed Redis Database instance.

## Attributes Reference

* `id` - The Managed Redis Database ID.

* `name` - The Managed Redis Database name.

* `cluster_id` - The Managed Redis Cluster ID that is hosting the Managed Redis Database.

* `linked_database_id` - The Linked Database list for the Managed Redis Database instance.

* `linked_database_group_nickname` - The Linked Database Group Nickname for the Managed Redis Database instance.

* `primary_access_key` - The Primary Access Key for the Managed Redis Database instance.

* `secondary_access_key` - The Secondary Access Key for the Managed Redis Database instance.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Redis Database.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Cache` - 2025-04-01
