---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_logger"
description: |-
  Manages a Logger within an API Management Service.
---

# azurerm_api_management_logger

Manages a Logger within an API Management Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "other"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_logger" "example" {
  name                = "example-logger"
  api_management_name = azurerm_api_management.example.name
  resource_group_name = azurerm_resource_group.example.name
  resource_id         = azurerm_application_insights.example.id

  application_insights {
    instrumentation_key = azurerm_application_insights.example.instrumentation_key
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of this Logger, which must be unique within the API Management Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `application_insights` - (Optional) An `application_insights` block as documented below. Changing this forces a new resource to be created.

* `buffered` - (Optional) Specifies whether records should be buffered in the Logger prior to publishing. Defaults to `true`.

* `description` - (Optional) A description of this Logger.

* `eventhub` - (Optional) An `eventhub` block as documented below. Changing this forces a new resource to be created.

* `resource_id` - (Optional) The target resource id which will be linked in the API-Management portal page. Changing this forces a new resource to be created.

---

An `application_insights` block supports the following:

* `connection_string` - (Optional) The connection string of Application Insights.

* `instrumentation_key` - (Optional) The instrumentation key used to push data to Application Insights.

~> **Note:** Either `connection_string` or `instrumentation_key` have to be specified.

---

An `eventhub` block supports the following:

* `name` - (Required) The name of an EventHub.

* `connection_string` - (Optional) The connection string of an EventHub Namespace.

~> **Note:** At least one of `connection_string` or `endpoint_uri` must be specified

* `user_assigned_identity_client_id` - (Optional) The Client Id of the User Assigned Identity with the "Azure Event Hubs Data Sender" role to the target EventHub Namespace. Required when `endpoint_uri` is set. If not specified the System Assigned Identity will be used.

* `endpoint_uri` - (Optional) The endpoint address of an EventHub Namespace. Required when `client_id` is set.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Logger.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Logger.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Logger.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Logger.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Logger.

## Import

API Management Loggers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_logger.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiManagement/service/example-apim/loggers/example-logger
```
