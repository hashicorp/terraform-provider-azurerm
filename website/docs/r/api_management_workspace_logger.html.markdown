---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_logger"
description: |-
  Manages a Logger within an API Management Workspace.
---

# azurerm_api_management_workspace_logger

Manages a Logger within an API Management Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"

  sku_name = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
}

resource "azurerm_application_insights" "example" {
  name                = "example-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_api_management_workspace_logger" "example" {
  name                        = "example-logger"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  description                 = "Example logger for workspace"

  application_insights {
    instrumentation_key = azurerm_application_insights.example.instrumentation_key
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the API Management Workspace Logger. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace. Changing this forces a new resource to be created.

---

* `application_insights` - (Optional) Specifies the application insights of the API Management Workspace Logger. Changing this forces a new resource to be created.

* `eventhub` - (Optional) Specifies the eventhub of the API Management Workspace Logger. Changing this forces a new resource to be created.

-> **Note:** Exactly one of `application_insights` or `eventhub` must be specified.

* `buffering_enabled` - (Optional) Specifies whether records should be buffered in the API Management Workspace Logger prior to publishing. Defaults to `true`.

* `description` - (Optional) Specifies a description of the API Management Workspace Logger.

* `resource_id` - (Optional) Specifies the target resource ID of the API Management Workspace Logger, which can be either an azure event hub or an application insights resource.

---

An `application_insights` block supports the following:

* `connection_string` - (Optional) Specifies the connection string of application insights.

* `instrumentation_key` - (Optional) Specifies the instrumentation key of the application insights.

-> **Note:** Exactly one of `connection_string` or `instrumentation_key` must be specified.

---

An `eventhub` block supports the following:

* `name` - (Required) Specifies the name of the eventhub.

* `connection_string` - (Optional) Specifies the connection string of the eventhub namespace.

* `endpoint_uri` - (Optional) Specifies the endpoint address of an eventhub namespace.

-> **Note:** Exactly one of `connection_string` or `endpoint_uri` must be specified.

* `user_assigned_identity_client_id` - (Optional) Specifies the client ID of user-assigned identity that has the "Azure Event Hubs Data Sender" role on the target eventhub namespace.

-> **Note:** * If this is omitted, the `SystemAssigned` will be used.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace Logger.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace Logger.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace Logger.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace Logger.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace Logger.

## Import

API Management Workspace Loggers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_logger.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/workspaces/workspace1/loggers/logger1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01