---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_service_azure_bot"
description: |-
  Manages an Azure Bot Service.
---

# azurerm_bot_service_azure_bot

Manages an Azure Bot Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "example" {
  name                    = "example-appinsightsapikey"
  application_insights_id = azurerm_application_insights.example.id
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

data "azurerm_client_config" "current" {}

resource "azurerm_bot_service_azure_bot" "example" {
  name                = "exampleazurebot"
  resource_group_name = azurerm_resource_group.example.name
  location            = "global"
  microsoft_app_id    = data.azurerm_client_config.current.client_id
  sku                 = "F0"

  endpoint                              = "https://example.com"
  developer_app_insights_api_key        = azurerm_application_insights_api_key.example.api_key
  developer_app_insights_application_id = azurerm_application_insights.example.app_id

  tags = {
    environment = "test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Bot Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Bot Service should exist. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the Azure Bot Service should exist. Changing this forces a new resource to be created.

* `microsoft_app_id` - (Required) The Microsoft Application ID for the Azure Bot Service. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU of the Azure Bot Service. Accepted values are `F0` or `S1`. Changing this forces a new resource to be created.

* `developer_app_insights_api_key` - (Optional) The Application Insights API Key to associate with this Azure Bot Service.

* `developer_app_insights_application_id` - (Optional) The resource ID of the Application Insights instance to associate with this Azure Bot Service.

* `developer_app_insights_key` - (Optional) The Application Insight Key to associate with this Azure Bot Service.

* `display_name` - (Optional) The name that the Azure Bot Service will be displayed as. This defaults to the value set for `name` if not specified.

* `endpoint` - (Optional) The Azure Bot Service endpoint.

* `microsoft_app_msi_id` - (Optional) The ID of the Microsoft App Managed Identity for this Azure Bot Service. Changing this forces a new resource to be created.

* `microsoft_app_tenant_id` - (Optional) The Tenant ID of the Microsoft App for this Azure Bot Service. Changing this forces a new resource to be created.

* `microsoft_app_type` - (Optional) The Microsoft App Type for this Azure Bot Service. Possible values are `MultiTenant`, `SingleTenant` and `UserAssignedMSI`. Changing this forces a new resource to be created.

* `local_authentication_enabled` - (Optional) Is local authentication enabled? Defaults to `true`.

* `luis_app_ids` - (Optional) A list of LUIS App IDs to associate with this Azure Bot Service.

* `luis_key` - (Optional) The LUIS key to associate with this Azure Bot Service.

* `streaming_endpoint_enabled` - (Optional) Is the streaming endpoint enabled for this Azure Bot Service. Defaults to `false`.

* `tags` - (Optional) A mapping of tags which should be assigned to this Azure Bot Service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Bot Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Bot Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Bot Service.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Bot Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Bot Service.

## Import

Azure Bot Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_service_azure_bot.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.BotService/botServices/botService1
```
