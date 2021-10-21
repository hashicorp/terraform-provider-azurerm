---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot"
description: |-
  Manages a Bot.
---

# azurerm_bot

Manages a Bot.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_bot" "example" {
  name                = "example"
  location            = "global"
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "F0"
  microsoft_app_id    = "${data.azurerm_client_config.current.client_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Bot. Changing this forces a new resource to be created. Must be globally unique.

* `resource_group_name` - (Required) The name of the resource group in which to create the Bot. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU of the Bot. Valid values include `F0` or `S1`. Changing this forces a new resource to be created.

* `microsoft_app_id` - (Required) The Microsoft Application ID for the Bot. Changing this forces a new resource to be created.

* `kind` - (Optional) The kind of Bot to create. This defaults to `azurebot` if not specified. Available values are `azurebot`, `bot`, `designer`, `function`, and `sdk`.

* `display_name` - (Optional) The name of the Bot will be displayed as. This defaults to `name` if not specified.

* `endpoint` - (Optional) The Bot endpoint.

* `developer_app_insights_key` - (Optional) The Application Insights Key to associate with the Bot.

* `developer_app_insights_api_key` - (Optional) The Application Insights API Key to associate with the Bot.

* `developer_app_insights_application_id` - (Optional) The Application Insights Application ID to associate with the Bot.

* `luis_app_ids` - (Optional) A list of LUIS App IDs to associate with the Bot.

* `luis_key` - (Optional) The LUIS key to associate with the Bot.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Bot.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Bot.
* `update` - (Defaults to 30 minutes) Used when updating the Bot.
* `read` - (Defaults to 5 minutes) Used when retrieving the Bot.
* `delete` - (Defaults to 30 minutes) Used when deleting the Bot.

## Import

Bot's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.BotService/botServices/example
```
