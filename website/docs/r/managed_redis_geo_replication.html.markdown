---
subcategory: "Managed Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_redis_geo_replication"
description: |-
  Manages Managed Redis Geo-Replication.
---

# azurerm_managed_redis_geo_replication

Manages Managed Redis Geo-Replication by linking and unlinking databases in a geo-replication group.

~> **Note:** This resource manages the geo-replication group membership for Managed Redis databases. All databases in the group must have the same `geo_replication_group_name` configured. Linking / unlinking databases will [discard all cache data and cause temporary instance outage](https://learn.microsoft.com/azure/redis/how-to-active-geo-replication).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-managedredis"
  location = "West Europe"
}

resource "azurerm_managed_redis" "amr1" {
  name                = "example-managedredis-amr1"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"

  sku_name = "Balanced_B3"

  default_database {
    geo_replication_group_name = "example-geo-group"
  }
}

resource "azurerm_managed_redis" "amr2" {
  name                = "example-managedredis-amr2"
  resource_group_name = azurerm_resource_group.example.name
  location            = "Central US"

  sku_name = "Balanced_B3"

  default_database {
    geo_replication_group_name = "example-geo-group"
  }
}

resource "azurerm_managed_redis_geo_replication" "example" {
  managed_redis_id = azurerm_managed_redis.amr1.id

  linked_managed_redis_ids = [
    azurerm_managed_redis.amr1.id,
    azurerm_managed_redis.amr2.id,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `managed_redis_id` - (Required) The ID of the Managed Redis through which geo-replication group will be managed. Linking is reciprocal, if A is linked to B, both A and B will have the same linking state. There is no need to have duplicate `azurerm_managed_redis_geo_replication` resources for each. Changing this forces a new resource to be created.

* `linked_managed_redis_ids` - (Required) A set of Managed Redis IDs to link together in the geo-replication group. Must include the `managed_redis_id` and can contain up to 5 Managed Redis IDs total. All Managed Redis must have the same `geo_replication_group_name` configured. Once linked, the geo-replication state of all databases will be updated.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Redis Database Geo-Replication resource (same as `managed_redis_id`).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Redis Database Geo-Replication.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Redis Database Geo-Replication.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Redis Database Geo-Replication.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Redis Database Geo-Replication.

## Import

Managed Redis Database Geo-Replication can be imported using the Managed Redis `resource id`, e.g.

```shell
terraform import azurerm_managed_redis_geo_replication.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redisEnterprise/cluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cache` - 2025-04-01
