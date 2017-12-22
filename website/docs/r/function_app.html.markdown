---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app"
sidebar_current: "docs-azurerm-resource-function-app-x"
description: |-
  Manages an Azure Functions service.

---

# azurerm_function_app

Manages an Azure Functions service.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "azure-functions-test-rg"
  location = "westus2"
}

resource "azurerm_storage_account" "test" {
  name                     = "azure-functions-test-sa"
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

* `name` - (Required) Specifies the name of the Azure Functions service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Azure Functions service.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Required) The ID of the App Service Plan within which to create this Azure Functions service. Changing this forces a new resource to be created.

* `storage_connection_string` - (Required) The connection string of the backend storage account which will be used by this Azure Functions service (such as the dashboard, logs).

* `app_settings` - (Optional) A key-value pair of App Settings.

* `enabled` - (Optional) Is the Azure Function service enabled? Changing this forces a new resource to be created.

* `version` - (Optional) The runtime version of this Azure Function service. Possible values are `~1` (this is the default value) and `beta`.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Azure Functions service

* `default_hostname` - The default hostname associated with the Azure Functions service - such as `mysite.azurewebsites.net`
