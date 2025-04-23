---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_application_live_view"
description: |-
  Manages a Spring Cloud Application Live View.
---

# azurerm_spring_cloud_application_live_view

-> **NOTE:** This resource is applicable only for Spring Cloud Service with enterprise tier.

Manages a Spring Cloud Application Live View.

!> Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_application_live_view` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_application_live_view" "example" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Application Live View. Changing this forces a new Spring Cloud Application Live View to be created. The only possible value is `default`.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new Spring Cloud Application Live View to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Application Live View.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application Live View.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application Live View.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application Live View.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application Live View.

## Import

Spring Cloud Application Live Views can be imported using the `resource id`, e.g.

```shellg
terraform import azurerm_spring_cloud_application_live_view.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/applicationLiveViews/default
```
