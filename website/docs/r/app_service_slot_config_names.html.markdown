---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_slot_config_names"
description: |-
  Manages Slot Configuration Names (within an App Service).
---

# azurerm_app_service_slot_config_names

Manages Slot Configuration Names (within an App Service).

-> **Note:** When using this resource - the `azurerm_app_service` resource or the `azurerm_app_service_slot` resource(s) should have `app_settings` keys or `connection_string` block defined, so that those settings can be sticked to the slot(s) after swapping. This resource is applied at app service site level.

## Example Usage

```hcl
resource "random_id" "server" {
  keepers = {
    azi_id = 1
  }

  byte_length = 8
}

resource "azurerm_resource_group" "example" {
  name     = "some-resource-group"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "example" {
  name                = "some-app-service-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = random_id.server.hex
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id

  app_settings = {
    "key1" = "value1"
    "key2" = "value2"
  }

  connection_string {
    name  = "Database1"
    type  = "SQLServer"
    value = "Server=some-server1.mydomain.com;Integrated Security=SSPI"
  }

  connection_string {
    name  = "Database2"
    type  = "SQLServer"
    value = "Server=some-server2.mydomain.com;Integrated Security=SSPI"
  }
}

resource "azurerm_app_service_slot" "example" {
  name                = random_id.server.hex
  app_service_name    = azurerm_app_service.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id

  app_settings = {
    "key1" = "value1"
    "key2" = "value2"
  }

  connection_string {
    name  = "Database1"
    type  = "SQLServer"
    value = "Server=some-server1.mydomain.com;Integrated Security=SSPI"
  }

  connection_string {
    name  = "Database2"
    type  = "SQLServer"
    value = "Server=some-server2.mydomain.com;Integrated Security=SSPI"
  }
}

resource "azurerm_app_service_slot_config_names" "example" {
  resource_group_name = "example"
  app_service_name    = "example"

  slot_config_names {
    app_setting_names       = ["key2"]
    connection_string_names = ["Database2"]
  }

  depends_on = [
    azurerm_app_service_slot.example
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `app_service_name` - (Required) The name of the App Service the slot configs to be applied.

* `resource_group_name` - (Required) The name of the Resource Group.

* `slot_config_names` - (Required) One or more `slot_config_names` blocks as defined below.

---

A `slot_config_names` block supports the following:

* `app_setting_names` - (Optional) Specifies a list of `app_settings` keys.

* `connection_string_names` - (Optional) Specifies a list of `connection_string` block.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The name of the App Service (Web Apps).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the App Service (Web Apps).
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service (Web Apps).
* `update` - (Defaults to 1 hour) Used when updating the App Service (Web Apps).
* `delete` - (Defaults to 1 hour) Used when deleting the App Service (Web Apps).

## Import

Slot Configuration Names does not support `terraform import`. It is Azure proxy only resource. This resource is not tracked by Azure Resource Manager.
