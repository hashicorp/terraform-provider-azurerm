---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_dynatrace_application_performance_monitoring"
description: |-
  Manages a Spring Cloud Application Performance Monitoring resource for Dynatrace.
---

# azurerm_spring_cloud_dynatrace_application_performance_monitoring

-> **Note:** This resource is only applicable for Spring Cloud Service enterprise tier

Manages a Spring Cloud Application Performance Monitoring resource for Dynatrace.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_dynatrace_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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

resource "azurerm_spring_cloud_dynatrace_application_performance_monitoring" "example" {
  name                    = "example"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
  globally_enabled        = true
  api_url                 = "https://example-api-url.com"
  api_token               = "dt0s01.AAAAAAAAAAAAAAAAAAAAAAAA.BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
  environment_id          = "example-environment-id"
  tenant                  = "example-tenant"
  tenant_token            = "dt0s01.AAAAAAAAAAAAAAAAAAAAAAAA.BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
  connection_point        = "https://example.live.dynatrace.com:443"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Application Performance Monitoring resource for Dynatrace. Changing this forces a new resource to be created.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new resource to be created.

* `tenant` - (Required) Specifies the Dynatrace tenant. 

* `tenant_token` - (Required) Specifies the internal token that is used for authentication when OneAgent connects to the Dynatrace cluster to send data.

* `connection_point` - (Required) Specifies the endpoint to connect to the Dynatrace environment.

* `globally_enabled` - (Optional) Specifies whether the Spring Cloud Application Performance Monitoring resource for Application Insights is enabled globally. Defaults to `false`.

* `api_url` - (Optional) Specifies the API Url of the Dynatrace environment.

* `api_token` - (Optional) Specifies the API token of the Dynatrace environment.

* `environment_id` - (Optional) Specifies the Dynatrace environment ID.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Spring Cloud Application Performance Monitoring resource for Dynatrace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application Performance Monitoring resource for Dynatrace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application Performance Monitoring resource for Dynatrace.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application Performance Monitoring resource for Dynatrace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application Performance Monitoring resource for Dynatrace.

## Import

Spring Cloud Application Performance Monitoring resource for Dynatrace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_dynatrace_application_performance_monitoring.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AppPlatform/spring/service1/apms/apm1
```
