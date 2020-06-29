---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_api_key"
description: |-
  Manages an Application Insights API key.
---

# azurerm_application_insights_api_key

Manages an Application Insights API key.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tf-test"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "tf-test-appinsights"
  location            = "West Europe"
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_application_insights_api_key" "read_telemetry" {
  name                    = "tf-test-appinsights-read-telemetry-api-key"
  application_insights_id = azurerm_application_insights.example.id
  read_permissions        = ["aggregate", "api", "draft", "extendqueries", "search"]
}

resource "azurerm_application_insights_api_key" "write_annotations" {
  name                    = "tf-test-appinsights-write-annotations-api-key"
  application_insights_id = azurerm_application_insights.example.id
  write_permissions       = ["annotations"]
}

resource "azurerm_application_insights_api_key" "authenticate_sdk_control_channel" {
  name                    = "tf-test-appinsights-authenticate-sdk-control-channel-api-key"
  application_insights_id = azurerm_application_insights.example.id
  read_permissions        = ["agentconfig"]
}

resource "azurerm_application_insights_api_key" "full_permissions" {
  name                    = "tf-test-appinsights-full-permissions-api-key"
  application_insights_id = azurerm_application_insights.example.id
  read_permissions        = ["agentconfig", "aggregate", "api", "draft", "extendqueries", "search"]
  write_permissions       = ["annotations"]
}

output "read_telemetry_api_key" {
  value = azurerm_application_insights_api_key.read_telemetry.api_key
}

output "write_annotations_api_key" {
  value = azurerm_application_insights_api_key.write_annotations.api_key
}

output "authenticate_sdk_control_channel" {
  value = azurerm_application_insights_api_key.authenticate_sdk_control_channel.api_key
}

output "full_permissions_api_key" {
  value = azurerm_application_insights_api_key.full_permissions.api_key
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Application Insights API key. Changing this forces a
    new resource to be created.

* `application_insights_id` - (Required) The ID of the Application Insights component on which the API key operates. Changing this forces a new resource to be created.

* `read_permissions` - (Optional) Specifies the list of read permissions granted to the API key. Valid values are `agentconfig`, `aggregate`, `api`, `draft`, `extendqueries`, `search`. Please note these values are case sensitive. Changing this forces a new resource to be created.

* `write_permissions` - (Optional) Specifies the list of write permissions granted to the API key. Valid values are `annotations`. Please note these values are case sensitive. Changing this forces a new resource to be created.

-> **Note:** At least one read or write permission must be defined.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Application Insights API key.

* `api_key` - The API Key secret (Sensitive).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Insights API Key.
* `update` - (Defaults to 30 minutes) Used when updating the Application Insights API Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Insights API Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Insights API Key.

## Import

Application Insights API keys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_insights_api_key.my_key /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.insights/components/instance1/apikeys/00000000-0000-0000-0000-000000000000
```

-> **Note:** The secret `api_key` cannot be retrieved during an import. You will need to edit the state by hand to set the secret value if you happen to have it backed up somewhere.
