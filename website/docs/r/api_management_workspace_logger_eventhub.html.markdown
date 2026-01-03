---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_logger_eventhub"
description: |-
  Manages an Event Hub Logger within an API Management Workspace.
---

# azurerm_api_management_workspace_logger_eventhub

Manages an Event Hub Logger within an API Management Workspace.

## Example Usage

### Using Connection String

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

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-eventhub-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "example" {
  name                = "example-eventhub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_api_management_workspace_logger_eventhub" "example" {
  name                        = "example-logger"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  description                 = "Example Event Hub logger for workspace"
  resource_id                 = azurerm_eventhub.example.id

  eventhub {
    name              = azurerm_eventhub.example.name
    connection_string = azurerm_eventhub_namespace.example.default_primary_connection_string
  }
}
```

### Using Managed Identity

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-identity"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"

  sku_name = "Premium_1"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id
    ]
  }
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-eventhub-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "example" {
  name                = "example-eventhub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_eventhub.example.id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_api_management_workspace_logger_eventhub" "example" {
  name                        = "example-logger"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  description                 = "Example Event Hub logger for workspace with managed identity"
  resource_id                 = azurerm_eventhub.example.id

  eventhub {
    name                             = azurerm_eventhub.example.name
    endpoint_uri                     = "${azurerm_eventhub_namespace.example.name}.servicebus.windows.net"
    user_assigned_identity_client_id = azurerm_user_assigned_identity.example.client_id
  }

  depends_on = [azurerm_role_assignment.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the API Management Workspace Logger. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace. Changing this forces a new resource to be created.

* `eventhub` - (Required) An `eventhub` block as defined below. Changing this forces a new resource to be created.

---

* `buffering_enabled` - (Optional) Specifies whether records should be buffered in the API Management Workspace Logger prior to publishing. Defaults to `true`.

* `description` - (Optional) Specifies a description of the API Management Workspace Logger.

* `resource_id` - (Optional) Specifies the Azure Resource ID of the Event Hub resource.

---

An `eventhub` block supports the following:

* `name` - (Required) Specifies the name of the Event Hub.

* `connection_string` - (Optional) Specifies the connection string of the Event Hub namespace.

* `endpoint_uri` - (Optional) Specifies the endpoint address of an Event Hub namespace.

-> **Note:** Exactly one of `connection_string` or `endpoint_uri` must be specified.

* `user_assigned_identity_client_id` - (Optional) Specifies the client ID of user-assigned identity that has the "Azure Event Hubs Data Sender" role on the target Event Hub namespace.

-> **Note:** If `endpoint_uri` is specified and this property is omitted, the `SystemAssigned` identity will be used.

-> **Note:** `user_assigned_identity_client_id` cannot be used with `connection_string`.

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

API Management Workspace Event Hub Loggers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_logger_eventhub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/workspaces/workspace1/loggers/logger1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01
