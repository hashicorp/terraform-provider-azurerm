---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_ms_teams"
description: |-
  Manages an MS Teams integration for a Bot Channel
---

# azurerm_bot_channel_ms_teams

Manages a MS Teams integration for a Bot Channel

~> **Note:** A bot can only have a single MS Teams Channel associated with it.

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
  microsoft_app_id    = data.azurerm_client_config.current.client_id
}

resource "azurerm_bot_channel_ms_teams" "example" {
  bot_name            = azurerm_bot_channels_registration.example.name
  location            = azurerm_bot_channels_registration.example.location
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which to create the Bot Channel. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `calling_web_hook` - (Optional) Specifies the webhook for Microsoft Teams channel calls.

* `deployment_environment` - (Optional) The deployment environment for Microsoft Teams channel calls. Possible values are `CommercialDeployment` and `GCCModerateDeployment`. Defaults to `CommercialDeployment`.

* `enable_calling` - (Optional) Specifies whether to enable Microsoft Teams channel calls. This defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Microsoft Teams Integration for a Bot Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Microsoft Teams Integration for a Bot Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft Teams Integration for a Bot Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Microsoft Teams Integration for a Bot Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Microsoft Teams Integration for a Bot Channel.

## Import

The Microsoft Teams Integration for a Bot Channel can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_ms_teams.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.BotService/botServices/example/channels/MsTeamsChannel
```
