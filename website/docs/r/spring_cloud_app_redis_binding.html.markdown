---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app_redis_binding"
description: |-
  Manages an Azure Spring Cloud Application Redis Binding instance.
---

# azurerm_spring_cloud_app_redis_binding

Manages an Azure Spring Cloud Application Redis Binding instance.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example-springcloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_spring_cloud_app" "example" {
  name                = "example-springcloudapp"
  resource_group_name = azurerm_resource_group.example.name
  service_name        = azurerm_spring_cloud_service.example.name
}

resource "azurerm_redis_cache" "example" {
  name                = "example-cache"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  capacity            = 0
  family              = "C"
  sku_name            = "Basic"
  enable_non_ssl_port = true
}

resource "azurerm_spring_cloud_app_redis_binding" "example" {
  name                = "example-bind"
  spring_cloud_app_id = azurerm_spring_cloud_app.example.id
  redis_cache_id      = azurerm_redis_cache.example.id
  redis_access_key    = azurerm_redis_cache.example.primary_access_key
  ssl_enabled         = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Application Binding. Changing this forces a new resource to be created.

* `spring_cloud_app_id` - (Required) Specifies the Spring Cloud Application resource ID in which the Binding is created. Changing this forces a new resource to be created.

* `redis_cache_id` - (Required) Specifies the Redis Cache resource ID. Changing this forces a new resource to be created.

* `redis_access_key` - (Required) Specifies the Redis Cache resource ID.

* `ssl_enabled` - (Optional) Is SSL used in the Spring Cloud Application Redis Binding? Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spring Cloud Application Redis Binding.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application Redis Binding.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application Redis Binding.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application Redis Binding.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application Redis Binding.

## Import

Spring Cloud Application Redis Binding can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_app_redis_binding.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.AppPlatform/Spring/myservice/apps/myapp/bindings/bind1
```
