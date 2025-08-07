---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_slack"
description: |-
  Manages a Slack integration for a Bot Channel
---

# azurerm_bot_channel_slack

Manages a Slack integration for a Bot Channel

~> **Note:** A bot can only have a single Slack Channel associated with it.

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

resource "azurerm_bot_channel_slack" "example" {
  bot_name            = azurerm_bot_channels_registration.example.name
  location            = azurerm_bot_channels_registration.example.location
  resource_group_name = azurerm_resource_group.example.name
  client_id           = "exampleId"
  client_secret       = "exampleSecret"
  verification_token  = "exampleVerificationToken"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which to create the Bot Channel. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `client_id` - (Required) The Client ID that will be used to authenticate with Slack.

* `client_secret` - (Required) The Client Secret that will be used to authenticate with Slack.

* `verification_token` - (Required) The Verification Token that will be used to authenticate with Slack.

* `landing_page_url` - (Optional) The Slack Landing Page URL.

* `signing_secret` - (Optional) The Signing Secret that will be used to sign the requests.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Slack Integration for a Bot Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Slack Integration for a Bot Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Slack Integration for a Bot Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Slack Integration for a Bot Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Slack Integration for a Bot Channel.

## Import

The Slack Integration for a Bot Channel can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_slack.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.BotService/botServices/example/channels/SlackChannel
```
