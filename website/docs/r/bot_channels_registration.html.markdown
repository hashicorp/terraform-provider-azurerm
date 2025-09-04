---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channels_registration"
description: |-
  Manages a Bot Channels Registration.
---

# azurerm_bot_channels_registration

Manages a Bot Channels Registration.

~> **Note:** Bot Channels Registration has been [deprecated by Azure](https://learn.microsoft.com/en-us/azure/bot-service/bot-service-resources-faq-azure?view=azure-bot-service-4.0#why-are-web-app-bot-and-bot-channel-registration-being-deprecated). New implementations should use the [`azurerm_bot_service_azure_bot`](./bot_service_azure_bot.html.markdown) resource.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Bot Channels Registration. Changing this forces a new resource to be created. Must be globally unique.

* `resource_group_name` - (Required) The name of the resource group in which to create the Bot Channels Registration. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU of the Bot Channels Registration. Valid values include `F0` or `S1`. Changing this forces a new resource to be created.

* `microsoft_app_id` - (Required) The Microsoft Application ID for the Bot Channels Registration. Changing this forces a new resource to be created.

* `cmk_key_vault_url` - (Optional) The CMK Key Vault Key URL to encrypt the Bot Channels Registration with the Customer Managed Encryption Key.

~> **Note:** It has to add the Key Vault Access Policy for the `Bot Service CMEK Prod` Service Principal and the `soft_delete_enabled` and the `purge_protection_enabled` is enabled on the `azurerm_key_vault` resource while using `cmk_key_vault_url`.

~> **Note:** It has to turn off the CMK feature before revoking Key Vault Access Policy. For more information, please refer to [Revoke access to customer-managed keys](https://docs.microsoft.com/azure/bot-service/bot-service-encryption?view=azure-bot-service-4.0&WT.mc_id=Portal-Microsoft_Azure_BotService#revoke-access-to-customer-managed-keys).

* `display_name` - (Optional) The name of the Bot Channels Registration will be displayed as. This defaults to `name` if not specified.

* `description` - (Optional) The description of the Bot Channels Registration.

* `endpoint` - (Optional) The Bot Channels Registration endpoint.

* `developer_app_insights_key` - (Optional) The Application Insights Key to associate with the Bot Channels Registration.

* `developer_app_insights_api_key` - (Optional) The Application Insights API Key to associate with the Bot Channels Registration.

* `developer_app_insights_application_id` - (Optional) The Application Insights Application ID to associate with the Bot Channels Registration.

* `icon_url` - (Optional) The icon URL to visually identify the Bot Channels Registration. Defaults to `https://docs.botframework.com/static/devportal/client/images/bot-framework-default.png`.

* `streaming_endpoint_enabled` - (Optional) Is the streaming endpoint enabled for the Bot Channels Registration. Defaults to `false`.

* `public_network_access_enabled` - (Optional) Is the Bot Channels Registration in an isolated network?

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Bot Channels Registration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Bot Channels Registration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Bot Channels Registration.
* `update` - (Defaults to 30 minutes) Used when updating the Bot Channels Registration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Bot Channels Registration.

## Import

Bot Channels Registration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channels_registration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.BotService/botServices/example
```
