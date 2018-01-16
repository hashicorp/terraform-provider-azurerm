---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app"
sidebar_current: "docs-azurerm-resource-function-app"
description: |-
  Manages a Function App.

---

# azurerm_function_app

Manages a Function App.

-> **Note:** Function Apps can be deployed to either an App Service Plan or to a Consumption Plan. At this time it's possible to deploy a Function App into an existing Consumption Plan or a new/existing App Service Plan [using the `azurerm_app_service_plan` Data Source](app_service_plan.html) - however it's not currently possible to create a new Consumption Plan natively in Terraform. Support for this will be added in the future, and in the interim can be achieved by using [the `azurerm_template_deployment` resource](template_deployment.html).

## Example Usage (with App Service Plan)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "azure-functions-test-rg"
  location = "westus2"
}

resource "azurerm_storage_account" "test" {
  name                     = "functionsapptestsa"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "azure-functions-test-service-plan"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "test-azure-functions"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Function App. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Function App.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Required) The ID of the App Service Plan within which to create this Function App. Changing this forces a new resource to be created.

* `storage_connection_string` - (Required) The connection string of the backend storage account which will be used by this Function App (such as the dashboard, logs).

* `app_settings` - (Optional) A key-value pair of App Settings.

* `enabled` - (Optional) Is the Function App enabled? Changing this forces a new resource to be created.

* `version` - (Optional) The runtime version associated with the Function App. Possible values are `~1` and `beta`. Defaults to `~1`.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Function App

* `default_hostname` - The default hostname associated with the Function App - such as `mysite.azurewebsites.net`


## Import

Function Apps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_function_app.functionapp1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/functionapp1
```
