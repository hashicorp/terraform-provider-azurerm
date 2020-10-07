---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_diagnostic"
description: |-
  Manages a API Management Service API Diagnostics Logs.
---

# azurerm_api_management_api_diagnostic

Manages a API Management Service API Diagnostics Logs.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_api" "example" {
  name                = "example-api"
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  revision            = "1"
  display_name        = "Example API"
  path                = "example"
  protocols           = ["https"]

  import {
    content_format = "swagger-link-json"
    content_value  = "http://conferenceapi.azurewebsites.net/?format=json"
  }
}

resource "azurerm_api_management_logger" "example" {
  name                = "example-apimlogger"
  api_management_name = azurerm_api_management.example.name
  resource_group_name = azurerm_resource_group.example.name

  application_insights {
    instrumentation_key = azurerm_application_insights.example.instrumentation_key
  }
}

resource "azurerm_api_management_api_diagnostic" "example" {
  resource_group_name      = azurerm_resource_group.example.name
  api_management_name      = azurerm_api_management.example.name
  api_name                 = azurerm_api_management_api.example.name
  api_management_logger_id = azurerm_api_management_logger.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_logger_id` - (Required) The ID (name) of the Diagnostics Logger.

* `api_management_name` - (Required) The name of the API Management Service instance. Changing this forces a new API Management Service API Diagnostics Logs to be created.

* `api_name` - (Required) The name of the API on which to configure the Diagnostics Logs. Changing this forces a new API Management Service API Diagnostics Logs to be created.

* `identifier` - (Required) Identifier of the Diagnostics Logs. Possible values are `applicationinsights` and `azuremonitor`. Changing this forces a new API Management Service API Diagnostics Logs to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the API Management Service API Diagnostics Logs should exist. Changing this forces a new API Management Service API Diagnostics Logs to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Service API Diagnostics Logs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Service API Diagnostics Logs.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Service API Diagnostics Logs.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Service API Diagnostics Logs.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Service API Diagnostics Logs.

## Import

API Management Service API Diagnostics Logs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_diagnostic.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/apis/api1/diagnostics/diagnostic1/loggers/logger1
```
