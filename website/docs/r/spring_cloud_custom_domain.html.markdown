---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_custom_domain"
description: |-
  Manages an Azure Spring Cloud Custom Domain.
---

# azurerm_spring_cloud_custom_domain

Manages an Azure Spring Cloud Custom Domain.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "example" {
  name = "example-resources"
}

data "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = data.azurerm_resource_group.example.name
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

resource "azurerm_dns_cname_record" "example" {
  name                = "record1"
  zone_name           = data.azurerm_dns_zone.example.name
  resource_group_name = data.azurerm_dns_zone.example.resource_group_name
  ttl                 = 300
  record              = azurerm_spring_cloud_app.example.fqdn
}

resource "azurerm_spring_cloud_custom_domain" "example" {
  name                = join(".", [azurerm_dns_cname_record.example.name, azurerm_dns_cname_record.example.zone_name])
  spring_cloud_app_id = azurerm_spring_cloud_app.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Custom Domain. Changing this forces a new resource to be created.

* `spring_cloud_app_id` - (Required) Specifies the resource ID of the Spring Cloud Application. Changing this forces a new resource to be created.

* `certificate_name` - (Optional) Specifies the name of the Spring Cloud Certificate that binds to the Spring Cloud Custom Domain.

* `thumbprint` - (Optional) Specifies the thumbprint of the Spring Cloud Certificate that binds to the Spring Cloud Custom Domain.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spring Cloud Custom Domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Custom Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Custom Domain.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Custom Domain.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Custom Domain.

## Import

Spring Cloud Custom Domain can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/Spring/spring1/apps/app1/domains/domain.com
```
