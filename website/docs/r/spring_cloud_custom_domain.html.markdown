---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_custom_domain"
description: |-
  Manages an Azure Spring Cloud Custom Domain.
---

# azurerm_spring_cloud_custom_domain

Manages an Azure Spring Cloud Custom Domain.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_custom_domain` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

data "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
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

* `certificate_name` - (Optional) Specifies the name of the Spring Cloud Certificate that binds to the Spring Cloud Custom Domain. Required when `thumbprint` is specified

* `thumbprint` - (Optional) Specifies the thumbprint of the Spring Cloud Certificate that binds to the Spring Cloud Custom Domain. Required when `certificate_name` is specified. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Custom Domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Custom Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Custom Domain.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Custom Domain.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Custom Domain.

## Import

Spring Cloud Custom Domain can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/apps/app1/domains/domain.com
```
