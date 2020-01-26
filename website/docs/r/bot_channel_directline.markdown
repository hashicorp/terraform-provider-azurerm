---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_directline"
description: |-
  Manages an Directline integration for a Bot Channel
---

# azurerm_bot_connection

Manages a Directline integration for a Bot Channel

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

resource "azurerm_bot_channel_directline" "example" {
  bot_name            = "${azurerm_bot_channels_registration.example.name}"
  location            = "${azurerm_bot_channels_registration.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  site {
    site_name  = "default"
    is_enabled = true
  }
}
```

## Argument Reference

The following arguments are supported:

- `resource_group_name` - (Required) The name of the resource group in which to create the Bot Channel. Changing this forces a new resource to be created.

- `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

- `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

- `site` - (Optional) A site represents a client application that you want to connect to your bot. Multiple `site` blocks may be defined as below

A `site` block as the following properties:

- `site_name` - (Required) The name of the site

- `is_enabled` - (Optional) Enables/Disables this site. Enabled by default

- `is_v1_enabled` - (Optional) Enables v1 of the Directline protocol for this site. Enabled by default

- `is_v3_enabled` - (Optional) Enables v3 of the Directline protocol for this site. Enabled by default

- `is_secure_site_enabled` - (Optional) Enables additional security measures for this site, see [Enhanced Directline Authentication Features](https://blog.botframework.com/2018/09/25/enhanced-direct-line-authentication-features). Disabled by default.

- `trusted_origins` - (Optional) This field is required when `is_secure_site_enabled` is enabled. Determines which origins can establish a Directline conversation for this site.

## Attributes Reference

The following attributes are exported:

- `id` - The Bot Channel ID.

## Import

The Microsoft Teams Channel for a Bot can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_ms_teams.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.BotService/botServices/example/channels/MsTeamsChannel
```
