---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_application_performance_monitoring"
description: |-
  Manages a Spring Cloud Application Performance Monitoring.
---

# azurerm_spring_cloud_dev_tool_portal

-> **NOTE:** This resource is applicable only for Spring Cloud Service with enterprise tier.

Manages a Spring Cloud Application Performance Monitoring.

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

resource "azurerm_spring_cloud_application_performance_monitoring" "example" {
  name                    = "example"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  type                    = "ApplicationInsights"
  properties = {
    any-string    = "any-string"
    sampling-rate = "12.0"
  }
  secrets = {
    connection-string = "XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXXX;XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXXXXXXXX"
  }
  globally_enabled = true
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Application Performance Monitoring. Changing this forces a new resource to be created.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new resource to be created.

* `type` - (Required) The type of the Spring Cloud Application Performance Monitoring. 

---

* `globally_enabled` - (Optional) Specifies whether the Spring Cloud Application Performance Monitoring is enabled globally. Defaults to `false`.

* `properties` - (Optional) Specifies a map of non-sensitive properties for launchProperties.

* `secrets` - (Optional) Specifies a map of sensitive properties for launchProperties.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Spring Cloud Application Performance Monitoring.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application Performance Monitoring.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application Performance Monitoring.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application Performance Monitoring.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application Performance Monitoring.

## Import

Spring Cloud Dev Tool Portals can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_application_performance_monitoring.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AppPlatform/spring/service1/apms/apm1
```
