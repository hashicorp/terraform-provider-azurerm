---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_connection"
description: |-
  Manages a Bot Connection.
---

# azurerm_bot_connection

Manages a Bot Connection.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_bot_channels_registration" "example" {
  name                = "example"
  location            = "global"
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "F0"
  microsoft_app_id    = "${data.azurerm_client_config.current.client_id}"
}

resource "azurerm_bot_connection" "example" {
  name                  = "example"
  bot_name              = azurerm_bot_channels_registration.example.name
  location              = azurerm_bot_channels_registration.example.location
  resource_group_name   = azurerm_resource_group.example.name
  service_provider_name = "box"
  client_id             = "exampleId"
  client_secret         = "exampleSecret"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Bot Connection. Changing this forces a new resource to be created. Must be globally unique.

* `resource_group_name` - (Required) The name of the resource group in which to create the Bot Connection. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this connection will be associated with. Changing this forces a new resource to be created.

* `service_provider_name` - (Required) The name of the service provider that will be associated with this connection. Changing this forces a new resource to be created.

* `client_id` - (Required) The Client ID that will be used to authenticate with the service provider.

* `client_secret` - (Required) The Client Secret that will be used to authenticate with the service provider.

* `scopes` - (Optional) The Scopes at which the connection should be applied.

* `parameters` - (Optional) A map of additional parameters to apply to the connection.

* `tags` - (Optional) A mapping of tags to assign to the resource.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Bot Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Bot Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Bot Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Bot Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Bot Connection.

## Import

Bot Connection can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.BotService/botServices/example/connections/example
```
