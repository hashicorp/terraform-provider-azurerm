---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_facebook"
description: |-
  Manages a Facebook integration for a Bot Channel
---

# azurerm_bot_channel_facebook

Manages a Facebook integration for a Bot Channel

~> **Note:** A bot can only have a single Facebook Channel associated with it.

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

resource "azurerm_bot_channel_facebook" "example" {
  bot_name                    = azurerm_bot_channels_registration.example.name
  location                    = azurerm_bot_channels_registration.example.location
  resource_group_name         = azurerm_resource_group.example.name
  facebook_application_id     = "563490254873576"
  facebook_application_secret = "8976d2536445ad5b976dee8437b9beb0"

  page {
    id           = "876248795081953"
    access_token = "CGGCec3UAFPMBAKwK3Ft8SEpO8ZCuvpNBI5DClaJCDfqJj2BgEHCKxcY0FDarmUQap6XxpZC9GWCW4nZCzjcKosAZAP7SO44X8Q8gAntbDIXgYUBGp9xtS8wUkwgKPobUePcOOVFkvClxvYZByuiQxoTiK9fQ9jZCPEorbmZCsKDZAx4VLnrNwCTZAPUwXxO61gfq4ZD"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group where the Facebook Channel should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `facebook_application_id` - (Required) The Facebook Application ID for the Facebook Channel.

* `facebook_application_secret` - (Required) The Facebook Application Secret for the Facebook Channel.

* `page` - (Required) One or more `page` blocks as defined below.

---

The `page` block supports the following:

* `id` - (Required) The Facebook Page ID for the Facebook Channel.

* `access_token` - (Required) The Facebook Page Access Token for the Facebook Channel.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Facebook Integration for a Bot Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Facebook Integration for a Bot Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Facebook Integration for a Bot Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Facebook Integration for a Bot Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Facebook Integration for a Bot Channel.

## Import

The Facebook Integration for a Bot Channel can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_facebook.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.BotService/botServices/botService1/channels/FacebookChannel
```
