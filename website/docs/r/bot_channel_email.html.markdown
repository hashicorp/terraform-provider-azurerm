---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_email"
description: |-
  Manages a Email integration for a Bot Channel
---

# azurerm_bot_channel_email

Manages a Email integration for a Bot Channel

~> **Note:** A bot can only have a single Email Channel associated with it.

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

resource "azurerm_bot_channel_email" "example" {
  bot_name            = azurerm_bot_channels_registration.example.name
  location            = azurerm_bot_channels_registration.example.location
  resource_group_name = azurerm_resource_group.example.name
  email_address       = "example.com"
  email_password      = "123456"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which to create the Bot Channel. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `email_address` - (Required) The email address that the Bot will authenticate with.

* `email_password` - (Optional) The email password that the Bot will authenticate with.

* `magic_code` - (Optional) The magic code used to set up OAUTH authentication.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Email Integration for a Bot Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Email Integration for a Bot Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Email Integration for a Bot Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Email Integration for a Bot Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Email Integration for a Bot Channel.

## Import

The Email Integration for a Bot Channel can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_email.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.BotService/botServices/example/channels/EmailChannel
```
