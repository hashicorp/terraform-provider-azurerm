---
subcategory: "Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redis_cache_patch_schedule"
description: |-
  Manages a Redis Cache Patch Schedule

---

# azurerm_redis_cache_patch_schedule

Manages a Redis Cache Patch Schedule.

~> **NOTE:** It's possible to define Redis Patch Schedule both within [the `azurerm_redis_cache` resource](redis_cache.html) via the `patch_schedule` block and by using [the `azurerm_redis_cache_patch_schedule` resource](redis_cache_patch_schedule.html). However it's not possible to use both methods to manage Patch Schedule within a Redis, since there'll be conflicts.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_redis_cache" "example" {
  name                = "example-cache"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  capacity            = 2
  family              = "C"
  sku_name            = "Standard"
  enable_non_ssl_port = false
  minimum_tls_version = "1.2"

  redis_configuration {
  }
}

resource azurerm_redis_cache_patch_schedule example {
  redis_cache_id = azurerm_redis_cache.example.id
  patch_schedule {
    day_of_week    = "Tuesday"
    start_hour_utc = 20
  }
}
```

## Argument Reference

The following arguments are supported:

* `redis_cache_id` - (Required) The ID of the Redis Cache. Changing this forces a new resource to be created.

* `patch_schedule` - (Required) A list of `patch_schedule` blocks as defined below.

---

A `patch_schedule` block supports the following:

* `day_of_week` - (Required) the Weekday name - possible values include `Monday`, `Tuesday`, `Wednesday` etc.

* `start_hour_utc` - (Optional) the Start Hour for maintenance in UTC - possible values range from `0 - 23`.

~> **Note:** The Patch Window lasts for `5` hours from the `start_hour_utc`.

* `maintenance_window` - (Optional) The ISO 8601 timespan which specifies the amount of time the Redis Cache can be updated. Defaults to `PT5H`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Redis Cache.

## Timeouts

 The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Redis Cache Patch Schedule.
* `update` - (Defaults to 30 minutes) Used when updating the Redis Cache Patch Schedule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Redis Cache Patch Schedule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Redis Cache Patch Schedule.

## Import

Redis Cache Patch Schedule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redis_cache_patch_schedule.patch1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redis/cache1
```
