---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app"
sidebar_current: "docs-azurerm-resource-app-service-function-app"
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

* `enable_builtin_logging` - (Optional) Should the built-in logging of this Function App be enabled? Defaults to `true`.

* `connection_string` - (Optional) An `connection_string` block as defined below.

* `client_affinity_enabled` - (Optional) Should the Function App send session affinity cookies, which route client requests in the same session to the same instance?

* `enabled` - (Optional) Is the Function App enabled?

* `https_only` - (Optional) Can the Function App only be accessed via HTTPS? Defaults to `false`.

* `version` - (Optional) The runtime version associated with the Function App. Defaults to `~1`.

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

* `always_on` - (Optional) Should the app be loaded at all times? Defaults to `false`.
* `default_documents` - (Optional) The ordering of default documents to load, if an address isn't specified.
* `dotnet_framework_version` - (Optional) The version of the .net framework's CLR used in this App Service. Possible values are `v2.0` (which will use the latest version of the .net framework for the .net CLR v2 - currently `.net 3.5`) and `v4.0` (which corresponds to the latest version of the .net CLR v4 - which at the time of writing is `.net 4.7.1`). [For more information on which .net CLR version to use based on the .net framework you're targeting - please see this table](https://en.wikipedia.org/wiki/.NET_Framework_version_history#Overview). Defaults to `v4.0`.
* `http2_enabled` - (Optional) Is HTTP2 Enabled on this App Service? Defaults to `false`.
* `ftps_state` - (Optional) State of FTP / FTPS service for this AppService. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`.
* `ip_restriction` - (Optional) One or more `ip_restriction` blocks as defined below.
* `java_version` - (Optional) The version of Java to use. If specified `java_container` and `java_container_version` must also be specified. Possible values are `1.7` and `1.8`.
* `java_container` - (Optional) The Java Container to use. If specified `java_version` and `java_container_version` must also be specified. Possible values are `JETTY` and `TOMCAT`.
* `java_container_version` - (Optional) The version of the Java Container to use. If specified `java_version` and `java_container` must also be specified.

* `local_mysql_enabled` - (Optional) Is "MySQL In App" Enabled? This runs a local MySQL instance with your app and shares resources from the App Service plan.

~> **NOTE:** MySQL In App is not intended for production environments and will not scale beyond a single instance. Instead you may wish [to use Azure Database for MySQL](/docs/providers/azurerm/r/mysql_database.html).

* `linux_fx_version` - (Optional) Linux App Framework and version for the AppService, e.g. `DOCKER|(golang:latest)`.
* `managed_pipeline_mode` - (Optional) The Managed Pipeline Mode. Possible values are `Integrated` and `Classic`. Defaults to `Integrated`.
* `min_tls_version` - (Optional) The minimum supported TLS version for the app service. Possible values are `1.0`, `1.1`, and `1.2`. Defaults to `1.2` for new app services.
* `php_version` - (Optional) The version of PHP to use in this App Service. Possible values are `5.5`, `5.6`, `7.0` and `7.1`.
* `python_version` - (Optional) The version of Python to use in this App Service. Possible values are `2.7` and `3.4`.
* `remote_debugging_enabled` - (Optional) Is Remote Debugging Enabled? Defaults to `false`.
* `remote_debugging_version` - (Optional) Which version of Visual Studio should the Remote Debugger be compatible with? Possible values are `VS2012`, `VS2013`, `VS2015` and `VS2017`.
* `scm_type` - (Optional) The type of Source Control enabled for this App Service. Possible values include `None` and `LocalGit`. Defaults to `None`.

~> **NOTE:** Additional Source Control types will be added in the future, once support for them has been added in the Azure SDK for Go.

* `use_32_bit_worker_process` - (Optional) Should the App Service run in 32 bit mode, rather than 64 bit mode?

~> **NOTE:** when using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `virtual_network_name` - (Optional) The name of the Virtual Network which this App Service should be attached to.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

---

`ip_restriction` supports the following:

* `ip_address` - (Required) The IP Address used for this IP Restriction.

* `subnet_mask` - (Optional) The Subnet mask used for this IP Restriction. Defaults to `255.255.255.255`.

---

`identity` supports the following:

* `type` - (Required) Specifies the identity type of the App Service. At this time the only allowed value is `SystemAssigned`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Function App

* `default_hostname` - The default hostname associated with the Function App - such as `mysite.azurewebsites.net`

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12`

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this App Service.

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this App Service.

---

`identity` exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this App Service.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this App Service.


`site_credential` exports the following:

* `username` - The username which can be used to publish to this App Service
* `password` - The password associated with the username, which can be used to publish to this App Service.


## Import

Function Apps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_function_app.functionapp1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/functionapp1
```
