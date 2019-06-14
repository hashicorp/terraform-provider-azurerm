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

* `app_command_line` - (Optional) App command line to launch, e.g. `/sbin/myserver -b 0.0.0.0`.

* `cors` - (Optional) A `cors` block as defined below.

* `default_documents` - (Optional) The ordering of default documents to load, if an address isn't specified.

* `dotnet_framework_version` - (Optional) The version of the .net framework's CLR used in this App Service. Possible values are `v2.0` (which will use the latest version of the .net framework for the .net CLR v2 - currently `.net 3.5`) and `v4.0` (which corresponds to the latest version of the .net CLR v4 - which at the time of writing is `.net 4.7.1`). [For more information on which .net CLR version to use based on the .net framework you're targeting - please see this table](https://en.wikipedia.org/wiki/.NET_Framework_version_history#Overview). Defaults to `v4.0`.

* `ftps_state` - (Optional) State of FTP / FTPS service for this App Service. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`.

* `http2_enabled` - (Optional) Is HTTP2 Enabled on this App Service? Defaults to `false`.

* `ip_restriction` - (Optional) A [List of objects](/docs/configuration/attr-as-blocks.html) representing ip restrictions as defined below.

* `java_version` - (Optional) The version of Java to use. If specified `java_container` and `java_container_version` must also be specified. Possible values are `1.7`, `1.8` and `11`.

* `java_container` - (Optional) The Java Container to use. If specified `java_version` and `java_container_version` must also be specified. Possible values are `JETTY` and `TOMCAT`.

* `java_container_version` - (Optional) The version of the Java Container to use. If specified `java_version` and `java_container` must also be specified.

* `local_mysql_enabled` - (Optional) Is "MySQL In App" Enabled? This runs a local MySQL instance with your app and shares resources from the App Service plan.

~> **NOTE:** MySQL In App is not intended for production environments and will not scale beyond a single instance. Instead you may wish [to use Azure Database for MySQL](/docs/providers/azurerm/r/mysql_database.html).

* `linux_fx_version` - (Optional) Linux App Framework and version for the App Service. Possible options are a Docker container (`DOCKER|<user/image:tag>`), a base-64 encoded Docker Compose file (`COMPOSE|${filebase64("compose.yml")}`) or a base-64 encoded Kubernetes Manifest (`KUBE|${filebase64("kubernetes.yml")}`).

* `windows_fx_version` - (Optional) The Windows Docker container image (`DOCKER|<user/image:tag>`)

Additional examples of how to run Containers via the `azurerm_app_service` resource can be found in [the `./examples/app-service` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/app-service).

* `managed_pipeline_mode` - (Optional) The Managed Pipeline Mode. Possible values are `Integrated` and `Classic`. Defaults to `Integrated`.

* `min_tls_version` - (Optional) The minimum supported TLS version for the app service. Possible values are `1.0`, `1.1`, and `1.2`. Defaults to `1.2` for new app services.

* `php_version` - (Optional) The version of PHP to use in this App Service. Possible values are `5.5`, `5.6`, `7.0`, `7.1` and `7.2`.

* `python_version` - (Optional) The version of Python to use in this App Service. Possible values are `2.7` and `3.4`.

* `remote_debugging_enabled` - (Optional) Is Remote Debugging Enabled? Defaults to `false`.

* `remote_debugging_version` - (Optional) Which version of Visual Studio should the Remote Debugger be compatible with? Possible values are `VS2012`, `VS2013`, `VS2015` and `VS2017`.

* `scm_type` - (Optional) The type of Source Control enabled for this App Service. Defaults to `None`. Possible values are: `BitbucketGit`, `BitbucketHg`, `CodePlexGit`, `CodePlexHg`, `Dropbox`, `ExternalGit`, `ExternalHg`, `GitHub`, `LocalGit`, `None`, `OneDrive`, `Tfs`, `VSO` and `VSTSRM`

* `use_32_bit_worker_process` - (Optional) Should the App Service run in 32 bit mode, rather than 64 bit mode?

~> **NOTE:** when using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `virtual_network_name` - (Optional) The name of the Virtual Network which this App Service should be attached to.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

---

A `cors` block supports the following:

* `allowed_origins` - (Optional) A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.

* `support_credentials` - (Optional) Are credentials supported?

---

`identity` supports the following:

* `type` - (Required) Specifies the identity type of the App Service. At this time the only allowed value is `SystemAssigned`.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Function App

* `default_hostname` - The default hostname associated with the Function App - such as `mysite.azurewebsites.net`

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12`

* `possible_outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12,52.143.43.17` - not all of which are necessarily in use. Superset of `outbound_ip_addresses`.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this App Service.

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this App Service.

* `kind` - The Function App kind - such as `functionapp,linux,container`

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
