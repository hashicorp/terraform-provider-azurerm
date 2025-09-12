---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app_redis_association"
description: |-
  Associates a [Spring Cloud Application](spring_cloud_app.html) with a [Redis Cache](redis_cache.html).
---

# azurerm_spring_cloud_app_redis_association

Associates a [Spring Cloud Application](spring_cloud_app.html) with a [Redis Cache](redis_cache.html).

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_app_redis_association` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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

resource "azurerm_spring_cloud_app_redis_association" "example" {
  name                = "example-bind"
  spring_cloud_app_id = azurerm_spring_cloud_app.example.id
  redis_cache_id      = azurerm_redis_cache.example.id
  redis_access_key    = azurerm_redis_cache.example.primary_access_key
  ssl_enabled         = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Application Association. Changing this forces a new resource to be created.

* `spring_cloud_app_id` - (Required) Specifies the Spring Cloud Application resource ID in which the Association is created. Changing this forces a new resource to be created.

* `redis_cache_id` - (Required) Specifies the Redis Cache resource ID. Changing this forces a new resource to be created.

* `redis_access_key` - (Required) Specifies the Redis Cache access key.

* `ssl_enabled` - (Optional) Should SSL be used when connecting to Redis? Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Application Redis Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application Redis Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application Redis Association.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application Redis Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application Redis Association.

## Import

Spring Cloud Application Redis Association can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_app_redis_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.AppPlatform/spring/myservice/apps/myapp/bindings/bind1
```
