---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_line"
description: |-
  Manages a Line integration for a Bot Channel
---

# azurerm_bot_channel_line

Manages a Line integration for a Bot Channel

~> **Note:** A bot can only have a single Line Channel associated with it.

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

resource "azurerm_bot_channel_line" "example" {
  bot_name            = azurerm_bot_channels_registration.example.name
  location            = azurerm_bot_channels_registration.example.location
  resource_group_name = azurerm_resource_group.example.name

  line_channel {
    access_token = "asdfdsdfTYUIOIoj1231hkjhk"
    secret       = "aagfdgfd123567"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group where the Line Channel should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `line_channel` - (Required) One or more `line_channel` blocks as defined below.

---

The `line_channel` block supports the following:

* `access_token` - (Required) The access token which is used to call the Line Channel API.

* `secret` - (Required) The secret which is used to access the Line Channel.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Line Integration for a Bot Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Line Integration for a Bot Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Line Integration for a Bot Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Line Integration for a Bot Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Line Integration for a Bot Channel.

## Import

The Line Integration for a Bot Channel can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_line.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.BotService/botServices/botService1/channels/LineChannel
```
