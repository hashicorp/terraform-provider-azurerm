---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_direct_line_speech"
description: |-
  Manages a Direct Line Speech integration for a Bot Channel
---

# azurerm_bot_channel_direct_line_speech

Manages a Direct Line Speech integration for a Bot Channel

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                = "example-cogacct"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "SpeechServices"
  sku_name            = "S0"
}

resource "azurerm_bot_channels_registration" "example" {
  name                = "example-bcr"
  location            = "global"
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "F0"
  microsoft_app_id    = data.azurerm_client_config.current.client_id
}

resource "azurerm_bot_channel_direct_line_speech" "example" {
  bot_name                     = azurerm_bot_channels_registration.example.name
  location                     = azurerm_bot_channels_registration.example.location
  resource_group_name          = azurerm_resource_group.example.name
  cognitive_service_location   = azurerm_cognitive_account.example.location
  cognitive_service_access_key = azurerm_cognitive_account.example.primary_access_key
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group where the Direct Line Speech Channel should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `cognitive_account_id` - (Optional) The ID of the Cognitive Account this Bot Channel should be associated with.

* `cognitive_service_access_key` - (Required) The access key to access the Cognitive Service.

* `cognitive_service_location` - (Required) Specifies the supported Azure location where the Cognitive Service resource exists.

* `custom_speech_model_id` - (Optional) The custom speech model id for the Direct Line Speech Channel.

* `custom_voice_deployment_id` - (Optional) The custom voice deployment id for the Direct Line Speech Channel.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Direct Line Speech Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Direct Line Speech Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Direct Line Speech Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Direct Line Speech Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Direct Line Speech Channel.

## Import

Direct Line Speech Channels can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_direct_line_speech.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.BotService/botServices/botService1/channels/DirectLineSpeechChannel
```
