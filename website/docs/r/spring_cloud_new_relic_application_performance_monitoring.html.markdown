---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_new_relic_application_performance_monitoring"
description: |-
  Manages a Spring Cloud Application Performance Monitoring resource for New Relic.
---

# azurerm_spring_cloud_new_relic_application_performance_monitoring

-> **Note:** This resource is only applicable for Spring Cloud Service enterprise tier

Manages a Spring Cloud Application Performance Monitoring resource for New Relic.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_new_relic_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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

resource "azurerm_spring_cloud_new_relic_application_performance_monitoring" "example" {
  name                    = "example"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
  app_name                = "example-app-name"
  license_key             = "example-license-key"
  app_server_port         = 8080
  labels = {
    tagName1 = "tagValue1"
    tagName2 = "tagValue2"
  }
  globally_enabled = true
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Application Performance Monitoring resource for New Relic. Changing this forces a new resource to be created.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new resource to be created.

* `app_name` - (Required) Specifies the application name used to report data to New Relic.

* `license_key` - (Required) Specifies the license key associated with the New Relic account. This key binds your agent's data to your account in New Relic service.

---

* `agent_enabled` - (Optional) Specifies whether enable the agent. Defaults to `true`.

* `app_server_port` - (Optional) Specifies the port number to differentiate JVMs for the same app on the same machine.

* `audit_mode_enabled` - (Optional) Specifies whether enable plain text logging of all data sent to New Relic to the agent logfile. Defaults to `false`.

* `auto_app_naming_enabled` - (Optional) Specifies whether enable the reporting of data separately for each web app. Defaults to `false`.

* `auto_transaction_naming_enabled` - (Optional) Specifies whether enable the component-based transaction naming. Defaults to `true`.

* `custom_tracing_enabled` - (Optional) Specifies whether enable all instrumentation using an `@Trace` annotation. Disabling this causes `@Trace` annotations to be ignored. Defaults to `true`.

* `labels` - (Optional) Specifies a mapping of labels to be added to the New Relic application.

* `globally_enabled` - (Optional) Specifies whether the Spring Cloud Application Performance Monitoring resource for Application Insights is enabled globally. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Application Performance Monitoring resource for New Relic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application Performance Monitoring resource for New Relic.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application Performance Monitoring resource for New Relic.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application Performance Monitoring resource for New Relic.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application Performance Monitoring resource for New Relic.

## Import

Spring Cloud Application Performance Monitoring resource for New Relic can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_new_relic_application_performance_monitoring.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AppPlatform/spring/service1/apms/apm1
```
