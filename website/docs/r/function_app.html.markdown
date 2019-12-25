---
subcategory: "App Service (Web Apps)"
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
resource "azurerm_resource_group" "example" {
  name     = "azure-functions-test-rg"
  location = "westus2"
}

resource "azurerm_storage_account" "example" {
  name                     = "functionsapptestsa"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "example" {
  name                = "azure-functions-test-service-plan"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "example" {
  name                      = "test-azure-functions"
  location                  = "${azurerm_resource_group.example.location}"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.example.id}"
  storage_connection_string = "${azurerm_storage_account.example.primary_connection_string}"
}
```
## Example Usage (in a Consumption Plan)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "azure-functions-cptest-rg"
  location = "westus2"
}

resource "azurerm_storage_account" "example" {
  name                     = "functionsapptestsa"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "example" {
  name                = "azure-functions-test-service-plan"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "example" {
  name                      = "test-azure-functions"
  location                  = "${azurerm_resource_group.example.location}"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.example.id}"
  storage_connection_string = "${azurerm_storage_account.example.primary_connection_string}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Function App. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Function App.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Required) The ID of the App Service Plan within which to create this Function App.

* `storage_connection_string` - (Required) The connection string of the backend storage account which will be used by this Function App (such as the dashboard, logs).

* `app_settings` - (Optional) A key-value pair of App Settings.

~> **Note:** When integrating a `CI/CD pipeline` and expecting to run from a deployed package in `Azure` you must seed your `app settings` as part of terraform code for function app to be successfully deployed. `Important Default key pairs`: (`"WEBSITE_RUN_FROM_PACKAGE" = ""`, `"FUNCTIONS_WORKER_RUNTIME" = "node"` (or python, etc), `"WEBSITE_NODE_DEFAULT_VERSION" = "10.14.1"`, `"APPINSIGHTS_INSTRUMENTATIONKEY" = ""`).

~> **Note:**  When using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `auth_settings` - (Optional) A `auth_settings` block as defined below.

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

* `always_on` - (Optional) Should the Function App be loaded at all times? Defaults to `false`.
* `use_32_bit_worker_process` - (Optional) Should the Function App run in 32 bit mode, rather than 64 bit mode? Defaults to `true`.

~> **Note:** when using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

* `virtual_network_name` - (Optional) The name of the Virtual Network which this App Service should be attached to.

* `linux_fx_version` - (Optional) Linux App Framework and version for the AppService, e.g. `DOCKER|(golang:latest)`.

* `http2_enabled` - (Optional) Specifies whether or not the http2 protocol should be enabled. Defaults to `false`.

* `min_tls_version` - (Optional) The minimum supported TLS version for the function app. Possible values are `1.0`, `1.1`, and `1.2`. Defaults to `1.2` for new function apps.

* `ftps_state` - (Optional) State of FTP / FTPS service for this function app. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`.

* `cors` - (Optional) A `cors` block as defined below.

---

A `cors` block supports the following:

* `allowed_origins` - (Optional) A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.

* `support_credentials` - (Optional) Are credentials supported?

---

`identity` supports the following:

* `type` - (Required) Specifies the identity type of the App Service. At this time the only allowed value is `SystemAssigned`.

---

An `auth_settings` block supports the following:

* `enabled` - (Required) Is Authentication enabled?

* `active_directory` - (Optional) A `active_directory` block as defined below.

* `additional_login_params` - (Optional) Login parameters to send to the OpenID Connect authorization endpoint when a user logs in. Each parameter must be in the form "key=value".

* `allowed_external_redirect_urls` - (Optional) External URLs that can be redirected to as part of logging in or logging out of the app.

* `default_provider` - (Optional) The default provider to use when multiple providers have been set up. Possible values are `AzureActiveDirectory`, `Facebook`, `Google`, `MicrosoftAccount` and `Twitter`.

~> **NOTE:** When using multiple providers, the default provider must be set for settings like `unauthenticated_client_action` to work.

* `facebook` - (Optional) A `facebook` block as defined below.

* `google` - (Optional) A `google` block as defined below.

* `issuer` - (Optional) Issuer URI. When using Azure Active Directory, this value is the URI of the directory tenant, e.g. https://sts.windows.net/{tenant-guid}/.

* `microsoft` - (Optional) A `microsoft` block as defined below.

* `runtime_version` - (Optional) The runtime version of the Authentication/Authorization module.

* `token_refresh_extension_hours` - (Optional) The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to 72.

* `token_store_enabled` - (Optional) If enabled the module will durably store platform-specific security tokens that are obtained during login flows. Defaults to false.

* `twitter` - (Optional) A `twitter` block as defined below.

* `unauthenticated_client_action` - (Optional) The action to take when an unauthenticated client attempts to access the app. Possible values are `AllowAnonymous` and `RedirectToLoginPage`.

---

An `active_directory` block supports the following:

* `client_id` - (Required) The Client ID of this relying party application. Enables OpenIDConnection authentication with Azure Active Directory.

* `client_secret` - (Optional) The Client Secret of this relying party application. If no secret is provided, implicit flow will be used.

* `allowed_audiences` (Optional) Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

---

A `facebook` block supports the following:

* `app_id` - (Required) The App ID of the Facebook app used for login

* `app_secret` - (Required) The App Secret of the Facebook app used for Facebook Login.

* `oauth_scopes` (Optional) The OAuth 2.0 scopes that will be requested as part of Facebook Login authentication. https://developers.facebook.com/docs/facebook-login

---

A `google` block supports the following:

* `client_id` - (Required) The OpenID Connect Client ID for the Google web application.

* `client_secret` - (Required) The client secret associated with the Google web application.

* `oauth_scopes` (Optional) The OAuth 2.0 scopes that will be requested as part of Google Sign-In authentication. https://developers.google.com/identity/sign-in/web/

---

A `microsoft` block supports the following:

* `client_id` - (Required) The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret` - (Required) The OAuth 2.0 client secret that was created for the app used for authentication.

* `oauth_scopes` (Optional) The OAuth 2.0 scopes that will be requested as part of Microsoft Account authentication. https://msdn.microsoft.com/en-us/library/dn631845.aspx

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
