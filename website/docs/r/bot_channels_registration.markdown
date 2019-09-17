---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channels_registration"
sidebar_current: "docs-azurerm-resource-bot_channels_registration"
description: |-
  Manages a Bot Channels Registration.
---

# azurerm_bot_channels_registration

Manages a Bot Channels Registration.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "northeurope"
}

resource "azurerm_bot_channels_registration" "example" {
  name                = "example"
  location            = "global"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "F0"
  microsoft_app_id    = "${data.azurerm_client_config.current.service_principal_application_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Bot Channels Registration. Changing this forces a new resource to be created. Must be globally unique.

* `resource_group_name` - (Required) The name of the resource group in which to create the Bot Channels Registration. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU of the Bot Channels Registration. Valid values include `F0` or `S1`. Changing this forces a new resource to be created.

* `microsoft_app_id` - (Required) The Microsoft Application ID for the Bot Channels Registration. Changing this forces a new resource to be created.

* `display_name` - (Optional) The name of the Bot Channels Registration will be displayed as. This defaults to `name` if not specified. 

* `endpoint` - (Optional) The Bot Channels Registration endpoint.

* `developer_app_insights_key` - (Optional) The Application Insights Key to associate with the Bot Channels Registration.

* `developer_app_insights_api_key` - (Optional) The Application Insights API Key to associate with the Bot Channels Registration.

* `developer_app_insights_application_id` - (Optional) The Application Insights Application ID to associate with the Bot Channels Registration.

* `tags` - (Optional) A mapping of tags to assign to the resource.


## Attributes Reference

The following attributes are exported:

* `id` - The Bot Channels Registration ID.

## Import

Bot Channels Registration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channels_registration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.BotService/botServices/example
```
