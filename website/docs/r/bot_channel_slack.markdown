---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_slack"
sidebar_current: "docs-azurerm-resource-bot-channel-slack"
description: |-
  Manages a Slack integration for a Bot Channel
---

# azurerm_bot_connection

Manages a Slack integration for a Bot Channel

~> **Note** A bot can only have a single Slack Channel associated with it.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}

resource "azurerm_bot_channels_registration" "example" {
  name                = "example"
  location            = "global"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "F0"
  microsoft_app_id    = "${data.azurerm_client_config.current.service_principal_application_id}"
}

resource "azurerm_bot_channel_slack" "example" {
  bot_name              = "${azurerm_bot_channels_registration.example.name}"
  location              = "${azurerm_bot_channels_registration.example.location}"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  client_id             = "exampleId"
  client_secret         = "exampleSecret"
  verification_token    = "exampleVerificationToken"
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


## Attributes Reference

The following attributes are exported:

* `id` - The Bot Channel ID.

## Import

The Slack Channel for a Bot can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_slack.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.BotService/botServices/example/channels/SlackChannel
```
