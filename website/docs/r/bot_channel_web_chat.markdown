---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_web_chat"
description: |-
  Manages a Web Chat integration for a Bot Channel
---

# azurerm_bot_channel_web_chat

Manages a Web Chat integration for a Bot Channel

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
  microsoft_app_id    = data.azurerm_client_config.current.service_principal_application_id
}

resource "azurerm_bot_channel_web_chat" "example" {
  bot_name            = azurerm_bot_channels_registration.example.name
  location            = azurerm_bot_channels_registration.example.location
  resource_group_name = azurerm_resource_group.example.name

  sites {
    site_name = "TestSite"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group where the Web Chat Channel should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `site_names` - (Required) A list of Web Chat Site names.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Web Chat Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web Chat Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Web Chat Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web Chat Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Web Chat Channel.

## Import

Web Chat Channels can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_web_chat.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.BotService/botServices/botService1/channels/WebChatChannel
```
