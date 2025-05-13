---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app_function"
description: |-
  Manages a Function App Function.
---

# azurerm_function_app_function

Manages a Function App Function.

## Example Usage - Basic HTTP Trigger

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-group"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
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

  site_config {
    application_stack {
      python_version = "3.9"
    }
  }
}

resource "azurerm_function_app_function" "example" {
  name            = "example-function-app-function"
  function_app_id = azurerm_linux_function_app.example.id
  language        = "Python"
  test_data = jsonencode({
    "name" = "Azure"
  })
  config_json = jsonencode({
    "bindings" = [
      {
        "authLevel" = "function"
        "direction" = "in"
        "methods" = [
          "get",
          "post",
        ]
        "name" = "req"
        "type" = "httpTrigger"
      },
      {
        "direction" = "out"
        "name"      = "$return"
        "type"      = "http"
      },
    ]
  })
}
```

## Example Usage - HTTP Trigger with code upload

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-group"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "example" {
  name                = "example-service-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Windows"
  sku_name            = "S1"
}

resource "azurerm_windows_function_app" "example" {
  name                = "example-function-app"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_plan_id     = azurerm_service_plan.example.id

  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  site_config {
    application_stack {
      dotnet_version = "6"
    }
  }
}

resource "azurerm_function_app_function" "example" {
  name            = "example-function-app-function"
  function_app_id = azurerm_windows_function_app.example.id
  language        = "CSharp"

  file {
    name    = "run.csx"
    content = file("exampledata/run.csx")
  }

  test_data = jsonencode({
    "name" = "Azure"
  })

  config_json = jsonencode({
    "bindings" = [
      {
        "authLevel" = "function"
        "direction" = "in"
        "methods" = [
          "get",
          "post",
        ]
        "name" = "req"
        "type" = "httpTrigger"
      },
      {
        "direction" = "out"
        "name"      = "$return"
        "type"      = "http"
      },
    ]
  })
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the function. Changing this forces a new resource to be created.

* `function_app_id` - (Required) The ID of the Function App in which this function should reside. Changing this forces a new resource to be created.

* `config_json` - (Required) The config for this Function in JSON format.

---

* `enabled` - (Optional) Should this function be enabled. Defaults to `true`.

* `file` - (Optional) A `file` block as detailed below. Changing this forces a new resource to be created.

* `language` - (Optional) The language the Function is written in. Possible values are `CSharp`, `Custom`, `Java`, `Javascript`, `Python`, `PowerShell`, and `TypeScript`.

~> **Note:** when using `Custom` language, you must specify the code handler in the `host.json` file for your function. See the [official docs](https://docs.microsoft.com/azure/azure-functions/functions-custom-handlers#hostjson) for more information.

* `test_data` - (Optional) The test data for the function.

---

A `file` block supports the following:

* `name` - (Required) The filename of the file to be uploaded. Changing this forces a new resource to be created.

* `content` - (Required) The content of the file. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Function App Function

* `config_url` - The URL of the configuration JSON.

* `invocation_url` - The invocation URL.

* `script_root_path_url` - The Script root path URL.

* `script_url` - The script URL.

* `secrets_file_url` - The URL for the Secrets File.

* `test_data_url` - The Test data URL.

* `url` - The function URL.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Function App Function.
* `read` - (Defaults to 5 minutes) Used when retrieving the Function App Function.
* `update` - (Defaults to 30 minutes) Used when updating the Function App Function.
* `delete` - (Defaults to 5 minutes) Used when deleting the Function App Function.

## Import

a Function App Function can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_function_app_function.example "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/functions/function1"
```
