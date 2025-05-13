---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_sms"
description: |-
  Manages a SMS integration for a Bot Channel
---

# azurerm_bot_channel_sms

Manages a SMS integration for a Bot Channel

~> **Note:** A bot can only have a single SMS Channel associated with it.

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

resource "azurerm_bot_channel_sms" "example" {
  bot_name                        = azurerm_bot_channels_registration.example.name
  location                        = azurerm_bot_channels_registration.example.location
  resource_group_name             = azurerm_resource_group.example.name
  sms_channel_account_security_id = "BG61f7cf5157f439b084e98256409c2815"
  sms_channel_auth_token          = "jh8980432610052ed4e29565c5e232f"
  phone_number                    = "+12313803556"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group where the SMS Channel should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `phone_number` - (Required) The phone number for the SMS Channel.

* `sms_channel_account_security_id` - (Required) The account security identifier (SID) for the SMS Channel.

* `sms_channel_auth_token` - (Required) The authorization token for the SMS Channel.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SMS Integration for a Bot Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SMS Integration for a Bot Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the SMS Integration for a Bot Channel.
* `update` - (Defaults to 30 minutes) Used when updating the SMS Integration for a Bot Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the SMS Integration for a Bot Channel.

## Import

The SMS Integration for a Bot Channel can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_sms.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.BotService/botServices/botService1/channels/SmsChannel
```
