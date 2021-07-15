---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_skype"
description: |-
  Manages a Skype integration for a Bot Channel
---

# azurerm_bot_channel_skype

Manages a Skype integration for a Bot Channel

~> **Note** A bot can only have a single Skype Channel associated with it.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_bot_channels_registration" "example" {
  name                = "example-bcr"
  location            = "global"
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "F0"
  microsoft_app_id    = data.azurerm_client_config.current.client_id
}

resource "azurerm_bot_channel_skype" "example" {
  bot_name            = azurerm_bot_channels_registration.example.name
  location            = azurerm_bot_channels_registration.example.location
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group where the Skype Channel should be created. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `calling_web_hook` - (Optional) The webhook for Skype channel calls.

* `enable_calling` - (Optional) Is Skype channel calls enabled?

* `enable_groups` - (Optional) Is Skype channel groups enabled?

* `enable_media_cards` - (Optional) Is Skype channel media cards enabled?

* `enable_messaging` - (Optional) Is Skype channel messaging enabled?

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Skype Integration for a Bot Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Skype Integration for a Bot Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Skype Integration for a Bot Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Skype Integration for a Bot Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Skype Integration for a Bot Channel.

## Import

The Skype Integration for a Bot Channel can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_skype.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.BotService/botServices/botService1/channels/SkypeChannel
```
