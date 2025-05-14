---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_application_insights_application_performance_monitoring"
description: |-
  Manages a Spring Cloud Application Performance Monitoring resource for Application Insights.
---

# azurerm_spring_cloud_application_insights_application_performance_monitoring

-> **Note:** This resource is only applicable for Spring Cloud Service enterprise tier

Manages a Spring Cloud Application Performance Monitoring resource for Application Insights.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_application_insights_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_application_insights_application_performance_monitoring" "example" {
  name                         = "example"
  spring_cloud_service_id      = azurerm_spring_cloud_service.example.id
  connection_string            = azurerm_application_insights.example.instrumentation_key
  globally_enabled             = true
  role_name                    = "test-role"
  role_instance                = "test-instance"
  sampling_percentage          = 50
  sampling_requests_per_second = 10
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Application Performance Monitoring resource for Application Insights. Changing this forces a new resource to be created.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new resource to be created.

---

* `globally_enabled` - (Optional) Specifies whether the Spring Cloud Application Performance Monitoring resource for Application Insights is enabled globally. Defaults to `false`.

* `connection_string` - (Optional) The instrumentation key used to push data to Application Insights.

* `role_name` - (Optional) Specifies the cloud role name used to label the component on the application map.

* `role_instance` - (Optional) Specifies the cloud role instance.
 
* `sampling_percentage` - (Optional) Specifies the percentage for fixed-percentage sampling.

* `sampling_requests_per_second` - (Optional) Specifies the number of requests per second for the rate-limited sampling. 

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Spring Cloud Application Performance Monitoring resource for Application Insights.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application Performance Monitoring resource for Application Insights.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application Performance Monitoring resource for Application Insights.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application Performance Monitoring resource for Application Insights.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application Performance Monitoring resource for Application Insights.

## Import

Spring Cloud Application Performance Monitoring resource for Application Insights can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_application_insights_application_performance_monitoring.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AppPlatform/spring/service1/apms/apm1
```
