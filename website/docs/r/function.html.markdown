---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function"
description: |-
  Manages a Azure Function.
---

# azurerm_function

Manages a Azure Function.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "azure-functions-cptest-rg"
  location = "westus2"
}

resource "azurerm_storage_account" "example" {
  name                     = "functionsapptestsa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "example" {
  name                = "azure-functions-test-service-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "example" {
  name                      = "test-azure-function-app"
  location                  = azurerm_resource_group.example.location
  resource_group_name       = azurerm_resource_group.example.name
  app_service_plan_id       = azurerm_app_service_plan.example.id
  storage_connection_string = azurerm_storage_account.example.primary_connection_string
}

resource "azurerm_function" "example" {
  name                = "test-azure-function"
  app_name            = azurerm_function_app.example.name
  resource_group_name = azurerm_resource_group.example.location
  config              = "{\"bindings\": [{\"type\": \"http\",\"direction\": \"out\",\"name\": \"res\"},{\"name\": \"req\",\"authLevel\": \"anonymous\",\"methods\": [\"post\"],\"direction\": \"in\",\"type\": \"httpTrigger\"}]}"
  files = {
    "index.js" : "module.exports = function (context, req) { context.done(); };"
  }
  test_data = "{\"method\":\"post\",\"headers\":[{\"name\":\"content-type\",\"value\":\"application/json\"}],\"body\":\"{\\\"key\\\":\\\"value\\\"}\"}"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Function. Changing this forces a new Azure Function to be created.

* `app_name` - (Required) The name of the Azure Function App, where the Function should be created in. Changing this forces a new Azure Function to be created.

* `config` - (Required) The Azure Function configuration as JSON string (see https://github.com/Azure/azure-functions-host/wiki/function.json)

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Function App should exist. Changing this forces a new Azure Function to be created.

---

* `files` - (Optional) Specifies a list of files which together make up the Function app ("filename" = "filecontent" pairs)

* `test_data` - (Optional) Test data used when testing via the Azure Portal.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Function.

* `config_href` - The `function.json` configuration file URI.

* `href` - The function endpoint URI.

* `invoke_url_template` - The function endpoint URI including the route params.

* `language` - The automatically determined script language of the created Function.

* `script_href` - The Script URI.

* `script_root_href` - The Script root path URI.

* `secrets_file_href` - The Secrets file URI.

* `test_data_href` - The Test data URI.

* `type` - The type of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Function.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Function.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Function.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Function.

## Import

Azure Functions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_function.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/sites/site1/functions/function1
```
