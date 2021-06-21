---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_redis_cache"
description: |-
  Manages a API Management Redis Cache.
---

# azurerm_api_management_redis_cache

Manages a API Management Redis Cache.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_redis_cache" "example" {
  name                = "example-cache"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  capacity            = 1
  family              = "C"
  sku_name            = "Basic"
  enable_non_ssl_port = false
  minimum_tls_version = "1.2"

  redis_configuration {
  }
}

resource "azurerm_api_management_redis_cache" "example" {
  name              = "example-Redis-Cache"
  api_management_id = azurerm_api_management.example.id
  connection_string = azurerm_redis_cache.example.primary_connection_string
  description       = "Redis cache instances"
  redis_cache_id    = azurerm_redis_cache.example.id
  cache_location    = "East Us"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this API Management Redis Cache. Changing this forces a new API Management Redis Cache to be created.

* `api_management_id` - (Required) The resource ID of the Api Management Service from which to create this external cache. Changing this forces a new API Management Redis Cache to be created.

* `connection_string` - (Required) The connection string to the Cache for Redis.

---

* `description` - (Optional) The description of the API Management Redis Cache.

* `redis_cache_id` - (Optional) The resource ID of the Cache for Redis.

* `cache_location` - (Optional) The location where to use cache from. Possible values are `default` and valid Azure regions. Defaults to `default`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management Redis Cache.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Redis Cache.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Redis Cache.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Redis Cache.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Redis Cache.

## Import

API Management Redis Caches can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_redis_cache.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/caches/cache1
```
