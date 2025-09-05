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
  microsoft_app_id    = data.azurerm_client_config.current.client_id
}

resource "azurerm_bot_channel_web_chat" "example" {
  bot_name            = azurerm_bot_channels_registration.example.name
  location            = azurerm_bot_channels_registration.example.location
  resource_group_name = azurerm_resource_group.example.name

  site {
    name = "TestSite"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group where the Web Chat Channel should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `site` - (Optional) A site represents a client application that you want to connect to your bot. One or more `site` blocks as defined below.

---

A `site` block has the following properties:

* `name` - (Required) The name of the site.

* `user_upload_enabled` - (Optional) Is the user upload enabled for this site? Defaults to `true`.

* `endpoint_parameters_enabled` - (Optional) Is the endpoint parameters enabled for this site?

* `storage_enabled` - (Optional) Is the storage site enabled for detailed logging? Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Web Chat Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Web Chat Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Web Chat Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Web Chat Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Web Chat Channel.

## Import

Web Chat Channels can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_web_chat.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.BotService/botServices/botService1/channels/WebChatChannel
```
