---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app"
sidebar_current: "docs-azurerm-resource-function-app"
description: |-
  Manages a Function App.

---

# azurerm_function_app

Manages a Function App.

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
## Example Usage (in a Consumption Plan)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "azure-functions-cptest-rg"
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
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
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

* `connection_string` - (Optional) An `connection_string` block as defined below.

* `client_affinity_enabled` - (Optional) Should the Function App send session affinity cookies, which route client requests in the same session to the same instance?

* `enabled` - (Optional) Is the Function App enabled?

* `https_only` - (Optional) Can the Function App only be accessed via HTTPS? Defaults to `false`.

* `version` - (Optional) The runtime version associated with the Function App. Possible values are `~1` and `beta`. Defaults to `~1`.

* `site_config` - (Optional) A `site_config` object as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`connection_string` supports the following:

* `name` - (Required) The name of the Connection String.
* `type` - (Required) The type of the Connection String. Possible values are `APIHub`, `Custom`, `DocDb`, `EventHub`, `MySQL`, `NotificationHub`, `PostgreSQL`, `RedisCache`, `ServiceBus`, `SQLAzure` and  `SQLServer`.
* `value` - (Required) The value for the Connection String.

---

`site_config` supports the following:

* `always_on` - (Optional) Should the Function App be loaded at all times? Defaults to `false`.
* `use_32_bit_worker_process` - (Optional) Should the Function App run in 32 bit mode, rather than 64 bit mode? Defaults to `true`.

~> **Note:** when using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

---

`identity` supports the following:

* `type` - (Required) Specifies the identity type of the App Service. At this time the only allowed value is `SystemAssigned`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Function App

* `default_hostname` - The default hostname associated with the Function App - such as `mysite.azurewebsites.net`

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12`

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this App Service.

---

`identity` exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this App Service.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this App Service.


## Import

Function Apps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_function_app.functionapp1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/functionapp1
```
