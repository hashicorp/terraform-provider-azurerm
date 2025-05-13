---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_elastic_application_performance_monitoring"
description: |-
  Manages a Spring Cloud Application Performance Monitoring resource for Elastic.
---

# azurerm_spring_cloud_elastic_application_performance_monitoring

-> **Note:** This resource is only applicable for Spring Cloud Service enterprise tier

Manages a Spring Cloud Application Performance Monitoring resource for Elastic.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_elastic_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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

resource "azurerm_spring_cloud_elastic_application_performance_monitoring" "example" {
  name                    = "example"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
  globally_enabled        = true
  application_packages    = ["org.example", "org.another.example"]
  service_name            = "example-service-name"
  server_url              = "http://127.0.0.1:8200"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Application Performance Monitoring resource for Elastic. Changing this forces a new resource to be created.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new resource to be created.

* `application_packages` - (Required) Specifies a list of the packages which should be used to determine whether a stack trace frame is an in-app frame or a library frame. This is a comma separated list of package names.

* `service_name` - (Required) Specifies the service name which is used to keep all the errors and transactions of your service together and is the primary filter in the Elastic APM user interface.

* `server_url` - (Required) Specifies the server URL. The URL must be fully qualified, including protocol (http or https) and port.

* `globally_enabled` - (Optional) Specifies whether the Spring Cloud Application Performance Monitoring resource for Application Insights is enabled globally. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Application Performance Monitoring resource for Elastic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application Performance Monitoring resource for Elastic.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application Performance Monitoring resource for Elastic.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application Performance Monitoring resource for Elastic.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application Performance Monitoring resource for Elastic.

## Import

Spring Cloud Application Performance Monitoring resource for Elastic can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_elastic_application_performance_monitoring.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AppPlatform/spring/service1/apms/apm1
```
