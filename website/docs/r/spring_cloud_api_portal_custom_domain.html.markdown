---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_api_portal_custom_domain"
description: |-
  Manages a Spring Cloud API Portal Domain.
---

# azurerm_spring_cloud_api_portal_custom_domain

Manages a Spring Cloud API Portal Domain.

-> **Note:** This resource is applicable only for Spring Cloud Service with enterprise tier.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_api_portal_custom_domain` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_gateway" "example" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
}

resource "azurerm_spring_cloud_api_portal" "example" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
  gateway_ids             = [azurerm_spring_cloud_gateway.example.id]
}

resource "azurerm_spring_cloud_api_portal_custom_domain" "example" {
  name                       = "example.com"
  spring_cloud_api_portal_id = azurerm_spring_cloud_api_portal.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud API Portal Domain. Changing this forces a new Spring Cloud API Portal Domain to be created.

* `spring_cloud_api_portal_id` - (Required) The ID of the Spring Cloud API Portal. Changing this forces a new Spring Cloud API Portal Domain to be created.

---

* `thumbprint` - (Optional) Specifies the thumbprint of the Spring Cloud Certificate that binds to the Spring Cloud API Portal Domain.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud API Portal Domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud API Portal Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud API Portal Domain.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud API Portal Domain.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud API Portal Domain.

## Import

Spring Cloud API Portal Domains can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_api_portal_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/apiPortals/apiPortal1/domains/domain1
```
