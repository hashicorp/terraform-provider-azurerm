---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_static_web_app_function_app_registration"
description: |-
  Manages a Static Web App Function App Registration.
---

# azurerm_static_web_app

Manages an App Service Static Web App Function App Registration.

~> **Note:** This resource registers the specified Function App to the `Production` build of the Static Web App.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_static_web_app" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesstorageacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "example" {
  name                = "example-service-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_function_app" "example" {
  name                = "example-function-app"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_plan_id     = azurerm_service_plan.example.id

  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  site_config {}

  lifecycle {
    ignore_changes = [auth_settings_v2]
  }
}

resource "azurerm_static_web_app_function_app_registration" "example" {
  static_web_app_id = azurerm_static_web_app.example.id
  function_app_id   = azurerm_linux_function_app.example.id
}

```

## Argument Reference

The following arguments are supported:

* `static_web_app_id` (Required) - The ID of the Static Web App to register the Function App to as a backend. Changing this forces a new resource to be created. 

* `function_app_id` (Required) - The ID of a Linux or Windows Function App to connect to the Static Web App as a Backend. Changing this forces a new resource to be created. 

~> **Note:** Only one Function App can be connected to a Static Web App. Multiple Function Apps are not currently supported.

~> **Note:** Connecting a Function App resource to a Static Web App resource updates the Function App to use AuthV2 and configures the `azure_static_web_app_v2` which may need to be accounted for by the use of `ignore_changes` depending on the existing `auth_settings_v2` configuration of the target Function App.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Static Web App Function App Registration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Static Web App Function App Registration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Static Web App Function App Registration.

## Import

Static Web App Function App Registration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_static_web_app_function_app_registration.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Web/staticSites/my-static-site1/userProvidedFunctionApps/myFunctionApp
```
