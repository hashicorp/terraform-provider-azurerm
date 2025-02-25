---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_app_slot_hybrid_connection"
description: |-
  Manages a Web App Slot Hybrid Connection.
---

# azurerm_web_app_slot_hybrid_connection

Manages a Web App Slot Hybrid Connection.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_service_plan" "example" {
  name                = "example-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Windows"
  sku_name            = "S1"
}

resource "azurerm_relay_namespace" "example" {
  name                = "example-relay"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard"
}

resource "azurerm_relay_hybrid_connection" "example" {
  name                 = "examplerhc1"
  resource_group_name  = azurerm_resource_group.example.name
  relay_namespace_name = azurerm_relay_namespace.example.name
}

resource "azurerm_windows_web_app" "example" {
  name                = "example-web-app"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_plan_id     = azurerm_service_plan.example.id

  site_config {}
}

resource "azurerm_windows_web_app_slot" "example" {
  name       = "slot"
  web_app_id = azurerm_windows_web_app.example.id

  storage_account_name = azurerm_storage_account.example.name

  site_config {}
}

resource "azurerm_web_app_slot_hybrid_connection" "example" {
  name       = azurerm_windows_web_app_slot.example.name
  web_app_id = azurerm_windows_web_app.example.id
  relay_id   = azurerm_relay_hybrid_connection.example.id
  hostname   = "myhostname.example"
  port       = 8081
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name of the Web App Deployment Slot for this Hybrid Connection. Changing this forces a new resource to be created.

* `function_app_id` - (Required) The ID of the Web App for this Hybrid Connection. Changing this forces a new resource to be created.

* `relay_id` - (Required) The ID of the Relay Hybrid Connection to use. Changing this forces a new resource to be created.

* `hostname` - (Required) The hostname of the endpoint.

* `port` - (Required) The port to use for the endpoint

---

* `send_key_name` - (Optional) The name of the Relay key with `Send` permission to use. Defaults to `RootManageSharedAccessKey`

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Web App Hybrid Connection

* `namespace_name` - The name of the Relay Namespace.

* `relay_name` - The name of the Relay in use.

* `send_key_value` - The Primary Access Key for the `send_key_name`

* `service_bus_namespace` - The Service Bus Namespace.

* `service_bus_suffix` - The suffix for the endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web App Hybrid Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Web App Hybrid Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web App Hybrid Connection.
* `delete` - (Defaults to 5 minutes) Used when deleting the Web App Hybrid Connection.

## Import

Web App Slot Hybrid Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_app_slot_hybrid_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/sites/site1/slots/slot1/hybridConnectionNamespaces/relay1/relays/hybridConnection1
```
