---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app"
description: |-
  Manages a Function App.

---

# azurerm_function_app

Manages a Function App.

~> **Note:** To connect an Azure Function App and a subnet within the same region `azurerm_app_service_virtual_network_swift_connection` can be used.
For an example, check the `azurerm_app_service_virtual_network_swift_connection` documentation.

## Example Usage (with App Service Plan)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "azure-functions-test-rg"
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

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "example" {
  name                       = "test-azure-functions"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  app_service_plan_id        = azurerm_app_service_plan.example.id
  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key
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
  name                       = "test-azure-functions"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  app_service_plan_id        = azurerm_app_service_plan.example.id
  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key
}
```
## Example Usage (Linux)

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
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "example" {
  name                       = "test-azure-functions"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  app_service_plan_id        = azurerm_app_service_plan.example.id
  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key
  os_type                    = "linux"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Function App. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Function App.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Required) The ID of the App Service Plan within which to create this Function App.

* `app_settings` - (Optional) A map of key-value pairs for [App Settings](https://docs.microsoft.com/en-us/azure/azure-functions/functions-app-settings) and custom values.

* `auth_settings` - (Optional) A `auth_settings` block as defined below.

* `connection_string` - (Optional) An `connection_string` block as defined below.

* `client_affinity_enabled` - (Optional) Should the Function App send session affinity cookies, which route client requests in the same session to the same instance?

* `daily_memory_time_quota` - (Optional) The amount of memory in gigabyte-seconds that your application is allowed to consume per day. Setting this value only affects function apps under the consumption plan. Defaults to `0`.

* `enabled` - (Optional) Is the Function App enabled?

* `enable_builtin_logging` - (Optional) Should the built-in logging of this Function App be enabled? Defaults to `true`.

* `https_only` - (Optional) Can the Function App only be accessed via HTTPS? Defaults to `false`.

* `identity` - (Optional) An `identity` block as defined below.

* `os_type` - (Optional) A string indicating the Operating System type for this function app. 

~> **NOTE:** This value will be `linux` for Linux derivatives, or an empty string for Windows (default). When set to `linux` you must also set `azurerm_app_service_plan` arguments as `kind = "Linux"` and `reserved = true`

* `site_config` - (Optional) A `site_config` object as defined below.

* `source_control` - (Optional) A `source_control` block, as defined below.

* `storage_account_name` - (Required) The backend storage account name which will be used by this Function App (such as the dashboard, logs).

* `storage_account_access_key` - (Required) The access key which will be used to access the backend storage account for the Function App.

~> **Note:** When integrating a `CI/CD pipeline` and expecting to run from a deployed package in `Azure` you must seed your `app settings` as part of terraform code for function app to be successfully deployed. `Important Default key pairs`: (`"WEBSITE_RUN_FROM_PACKAGE" = ""`, `"FUNCTIONS_WORKER_RUNTIME" = "node"` (or python, etc), `"WEBSITE_NODE_DEFAULT_VERSION" = "10.14.1"`, `"APPINSIGHTS_INSTRUMENTATIONKEY" = ""`).

~> **Note:**  When using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `version` - (Optional) The runtime version associated with the Function App. Defaults to `~1`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`connection_string` supports the following:

* `name` - (Required) The name of the Connection String.

* `type` - (Required) The type of the Connection String. Possible values are `APIHub`, `Custom`, `DocDb`, `EventHub`, `MySQL`, `NotificationHub`, `PostgreSQL`, `RedisCache`, `ServiceBus`, `SQLAzure` and  `SQLServer`.

* `value` - (Required) The value for the Connection String.

---

`site_config` supports the following:

* `always_on` - (Optional) Should the Function App be loaded at all times? Defaults to `false`.

* `cors` - (Optional) A `cors` block as defined below.

* `ftps_state` - (Optional) State of FTP / FTPS service for this function app. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`. Defaults to `AllAllowed`.

* `health_check_path` - (Optional) Path which will be checked for this function app health.

* `http2_enabled` - (Optional) Specifies whether or not the http2 protocol should be enabled. Defaults to `false`.

* `ip_restriction` - (Optional) A [List of objects](/docs/configuration/attr-as-blocks.html) representing ip restrictions as defined below.

-> **NOTE** User has to explicitly set `ip_restriction` to empty slice (`[]`) to remove it.

* `java_version` - (Optional) Java version hosted by the function app in Azure. Possible values are `1.8`, `11`.

* `linux_fx_version` - (Optional) Linux App Framework and version for the AppService, e.g. `DOCKER|(golang:latest)`.

* `min_tls_version` - (Optional) The minimum supported TLS version for the function app. Possible values are `1.0`, `1.1`, and `1.2`. Defaults to `1.2` for new function apps.

* `pre_warmed_instance_count` - (Optional) The number of pre-warmed instances for this function app. Only affects apps on the Premium plan.

* `scm_ip_restriction` - (Optional) A [List of objects](/docs/configuration/attr-as-blocks.html) representing ip restrictions as defined below.

-> **NOTE** User has to explicitly set `scm_ip_restriction` to empty slice (`[]`) to remove it.

* `scm_type` - (Optional) The type of Source Control used by the Function App. Valid values include: `BitBucketGit`, `BitBucketHg`, `CodePlexGit`, `CodePlexHg`, `Dropbox`, `ExternalGit`, `ExternalHg`, `GitHub`, `LocalGit`, `None` (default), `OneDrive`, `Tfs`, `VSO`, and `VSTSRM`

~> **NOTE:** This setting is incompatible with the `source_control` block which updates this value based on the setting provided. 

* `scm_use_main_ip_restriction` - (Optional)  IP security restrictions for scm to use main. Defaults to false.  

-> **NOTE** Any `scm_ip_restriction` blocks configured are ignored by the service when `scm_use_main_ip_restriction` is set to `true`. Any scm restrictions will become active if this is subsequently set to `false` or removed.  

* `use_32_bit_worker_process` - (Optional) Should the Function App run in 32 bit mode, rather than 64 bit mode? Defaults to `true`.

~> **Note:** when using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

---

A `cors` block supports the following:

* `allowed_origins` - (Optional) A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.

* `support_credentials` - (Optional) Are credentials supported?

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Function App. Possible values are `SystemAssigned` (where Azure will generate a Service Principal for you), `UserAssigned` where you can specify the Service Principal IDs in the `identity_ids` field, and `SystemAssigned, UserAssigned` which assigns both a system managed identity as well as the specified user assigned identities.

~> **NOTE:** When `type` is set to `SystemAssigned`, The assigned `principal_id` and `tenant_id` can be retrieved after the Function App has been created. More details are available below.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned. Required if `type` is `UserAssigned`.

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

---

A `ip_restriction` block supports the following:

* `ip_address` - (Optional) The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

-> **NOTE:** One of either `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified

* `name` - (Optional) The name for this IP Restriction.

* `priority` - (Optional) The priority for this IP Restriction. Restrictions are enforced in priority order. By default, the priority is set to 65000 if not specified.

* `action` - (Optional) Does this restriction `Allow` or `Deny` access for this IP range. Defaults to `Allow`.  

---

A `scm_ip_restriction` block supports the following:

* `ip_address` - (Optional) The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

-> **NOTE:** One of either `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified

* `name` - (Optional) The name for this IP Restriction.

* `priority` - (Optional) The priority for this IP Restriction. Restrictions are enforced in priority order. By default, priority is set to 65000 if not specified.  

* `action` - (Optional) Allow or Deny access for this IP range. Defaults to Allow.  

---

A `source_control` block supports the following:

* `repo_url` - (Required) The URL of the source code repository.

* `branch` - (Optional) The branch of the remote repository to use. Defaults to 'master'. 

* `manual_integration` - (Optional) Limits to manual integration. Defaults to `false` if not specified. 

* `rollback_enabled` - (Optional) Enable roll-back for the repository. Defaults to `false` if not specified.

* `use_mercurial` - (Optional) Use Mercurial if `true`, otherwise uses Git. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Function App

* `custom_domain_verification_id` - An identifier used by App Service to perform domain ownership verification via DNS TXT record. 

* `default_hostname` - The default hostname associated with the Function App - such as `mysite.azurewebsites.net`

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12`

* `possible_outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12,52.143.43.17` - not all of which are necessarily in use. Superset of `outbound_ip_addresses`.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this App Service.

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this App Service.

* `kind` - The Function App kind - such as `functionapp,linux,container`

---

The `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this App Service.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this App Service.

---

The `site_credential` block exports the following:

* `username` - The username which can be used to publish to this App Service

* `password` - The password associated with the username, which can be used to publish to this App Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Function App.
* `update` - (Defaults to 30 minutes) Used when updating the Function App.
* `read` - (Defaults to 5 minutes) Used when retrieving the Function App.
* `delete` - (Defaults to 30 minutes) Used when deleting the Function App.

## Import

Function Apps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_function_app.functionapp1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/functionapp1
```
