---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_hybrid_connection"
description: |-
  Manages an App Service Hybrid Connection for an existing App Service, Relay and Service Bus.

---

# azurerm_app_service_hybrid_connection

Manages an App Service Hybrid Connection for an existing App Service, Relay and Service Bus.

## Example Usage

This example provisions an App Service, a Relay Hybrid Connection, and a Service Bus using their outputs to create the App Service Hybrid Connection.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exampleResourceGroup1"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "example" {
  name                = "exampleAppServicePlan1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "exampleAppService1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id
}

resource "azurerm_relay_namespace" "example" {
  name                = "exampleRN1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "Standard"
}

resource "azurerm_relay_hybrid_connection" "example" {
  name                 = "exampleRHC1"
  resource_group_name  = azurerm_resource_group.example.name
  relay_namespace_name = azurerm_relay_namespace.example.name
  user_metadata        = "examplemetadata"
}

resource "azurerm_app_service_hybrid_connection" "example" {
  app_service_name    = azurerm_app_service.example.name
  resource_group_name = azurerm_resource_group.example.name
  relay_id            = azurerm_relay_hybrid_connection.example.id
  hostname            = "testhostname.example"
  port                = 8080
  send_key_name       = "exampleSharedAccessKey"
}

```

## Argument Reference

The following arguments are supported:

* `app_service_name` - (Required) Specifies the name of the App Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the App Service.  Changing this forces a new resource to be created.

* `relay_id` - (Required) The ID of the Service Bus Relay. Changing this forces a new resource to be created.

* `hostname` - (Optional) The hostname of the endpoint.

* `port` - (Required) The port of the endpoint.

* `send_key_name` - (Optional) The name of the Service Bus key which has Send permissions. Defaults to `RootManageSharedAccessKey`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service.

* `namespace_name` - The name of the Relay Namespace.

* `send_key_value` - The value of the Service Bus Primary Access key.

* `service_bus_namespace` - The name of the Service Bus namespace.

* `service_bus_suffix` - The suffix for the service bus endpoint.

## Import

App Service Hybrid Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_hybrid_connection.example /subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/exampleResourceGroup1/providers/Microsoft.Web/sites/exampleAppService1/hybridConnectionNamespaces/exampleRN1/relays/exampleRHC1
```
