---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app_slot"
description: |-
  Manages a Function App Deployment Slot.

---

# azurerm_function_app_slot

Manages a Function App deployment Slot.

!> **Note:** This resource has been deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use [`azurerm_linux_function_app_slot`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/linux_function_app_slot) and [`azurerm_windows_function_app_slot`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/windows_function_app_slot) resources instead.

## Example Usage (with App Service Plan)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "azure-functions-test-rg"
  location = "West Europe"
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

resource "azurerm_function_app_slot" "example" {
  name                       = "test-azure-functions_slot"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  app_service_plan_id        = azurerm_app_service_plan.example.id
  function_app_name          = azurerm_function_app.example.name
  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Function App. Changing this forces a new resource to be created. 

* `resource_group_name` - (Required) The name of the resource group in which to create the Function App Slot. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Required) The ID of the App Service Plan within which to create this Function App Slot. Changing this forces a new resource to be created.

* `function_app_name` - (Required) The name of the Function App within which to create the Function App Slot. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) The backend storage account name which will be used by the Function App (such as the dashboard, logs). Changing this forces a new resource to be created.

* `storage_account_access_key` - (Required) The access key which will be used to access the backend storage account for the Function App.

* `app_settings` - (Optional) A key-value pair of App Settings.

~> **Note:** When integrating a `CI/CD pipeline` and expecting to run from a deployed package in `Azure` you must seed your `app settings` as part of terraform code for function app to be successfully deployed. `Important Default key pairs`: (`"WEBSITE_RUN_FROM_PACKAGE" = ""`, `"FUNCTIONS_WORKER_RUNTIME" = "node"` (or Python, etc), `"WEBSITE_NODE_DEFAULT_VERSION" = "10.14.1"`, `"APPINSIGHTS_INSTRUMENTATIONKEY" = ""`).

~> **Note:** The values for `AzureWebJobsStorage` and `FUNCTIONS_EXTENSION_VERSION` will be filled by other input arguments and shouldn't be configured separately. `AzureWebJobsStorage` is filled based on `storage_account_name` and `storage_account_access_key`. `FUNCTIONS_EXTENSION_VERSION` is filled based on `version`.

~> **Note:** When using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `auth_settings` - (Optional) An `auth_settings` block as defined below.

* `enable_builtin_logging` - (Optional) Should the built-in logging of the Function App be enabled? Defaults to `true`.

* `connection_string` - (Optional) A `connection_string` block as defined below.

* `os_type` - (Optional) A string indicating the Operating System type for this function app. The only possible value is `linux`. Changing this forces a new resource to be created.

~> **Note:** This value will be `linux` for Linux Derivatives or an empty string for Windows (default).

* `enabled` - (Optional) Is the Function App enabled? Defaults to `true`.

* `https_only` - (Optional) Can the Function App only be accessed via HTTPS? Defaults to `false`.

* `version` - (Optional) The runtime version associated with the Function App. Defaults to `~1`.

* `daily_memory_time_quota` - (Optional) The amount of memory in gigabyte-seconds that your application is allowed to consume per day. Setting this value only affects function apps under the consumption plan.

* `site_config` - (Optional) A `site_config` object as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `connection_string` block supports the following:

* `name` - (Required) The name of the Connection String.
* `type` - (Required) The type of the Connection String. Possible values are `APIHub`, `Custom`, `DocDb`, `EventHub`, `MySQL`, `NotificationHub`, `PostgreSQL`, `RedisCache`, `ServiceBus`, `SQLAzure` and `SQLServer`.
* `value` - (Required) The value for the Connection String.

---

The `site_config` block supports the following:

* `always_on` - (Optional) Should the Function App be loaded at all times? Defaults to `false`.

* `use_32_bit_worker_process` - (Optional) Should the Function App run in 32 bit mode, rather than 64 bit mode? Defaults to `true`.

~> **Note:** when using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

* `linux_fx_version` - (Optional) Linux App Framework and version for the AppService, e.g. `DOCKER|(golang:latest)`.

* `java_version` - (Optional) Java version hosted by the function app in Azure. Possible values are `1.8`, `11` & `17` (In-Preview).

* `http2_enabled` - (Optional) Specifies whether or not the HTTP2 protocol should be enabled. Defaults to `false`.

* `min_tls_version` - (Optional) The minimum supported TLS version for the function app. Possible values are `1.0`, `1.1`, and `1.2`. Defaults to `1.2` for new function apps.

* `ftps_state` - (Optional) State of FTP / FTPS service for this function app. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`.

* `health_check_path` - (Optional) Path which will be checked for this function app health.

* `app_scale_limit` - (Optional) The number of workers this function app can scale out to. Only applicable to apps on the Consumption and Premium plan.

* `runtime_scale_monitoring_enabled` - (Optional) Should Runtime Scale Monitoring be enabled?. Only applicable to apps on the Premium plan. Defaults to `false`.

* `scm_ip_restriction` - (Optional) A list of `scm_ip_restriction` objects representing IP restrictions as defined below.

-> **Note:** User has to explicitly set `scm_ip_restriction` to empty slice (`[]`) to remove it.

* `scm_type` - (Optional) The type of Source Control used by this function App. Valid values include: `BitBucketGit`, `BitBucketHg`, `CodePlexGit`, `CodePlexHg`, `Dropbox`, `ExternalGit`, `ExternalHg`, `GitHub`, `LocalGit`, `None` (default), `OneDrive`, `Tfs`, `VSO`, and `VSTSRM`.

~> **Note:** This setting is incompatible with the `source_control` block which updates this value based on the setting provided.

* `scm_use_main_ip_restriction` - (Optional) IP security restrictions for scm to use main. Defaults to `false`.

-> **Note:** Any `scm_ip_restriction` blocks configured are ignored by the service when `scm_use_main_ip_restriction` is set to `true`. Any scm restrictions will become active if this is subsequently set to `false` or removed.

* `elastic_instance_minimum` - (Optional) The number of minimum instances for this function app. Only applicable to apps on the Premium plan.

* `pre_warmed_instance_count` - (Optional) The number of pre-warmed instances for this function app. Only affects apps on the Premium plan.

* `dotnet_framework_version` - (Optional) The version of the .NET framework's CLR used in this function app. Possible values are `v4.0` (including .NET Core 2.1 and 3.1), `v5.0` and `v6.0`. [For more information on which .NET Framework version to use based on the runtime version you're targeting - please see this table](https://docs.microsoft.com/azure/azure-functions/functions-dotnet-class-library#supported-versions). Defaults to `v4.0`.

* `cors` - (Optional) A `cors` block as defined below.

* `ip_restriction` - (Optional) A list of `ip_restriction` objects representing IP restrictions as defined below.

* `auto_swap_slot_name` - (Optional) The name of the slot to automatically swap to during deployment

---

A `cors` block supports the following:

* `allowed_origins` - (Required) A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.

* `support_credentials` - (Optional) Are credentials supported?

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Function App. Possible values are `SystemAssigned` (where Azure will generate a Service Principal for you), `UserAssigned` where you can specify the Service Principal IDs in the `identity_ids` field, and `SystemAssigned, UserAssigned` which assigns both a system managed identity as well as the specified user assigned identities.

~> **Note:** When `type` is set to `SystemAssigned`, The assigned `principal_id` and `tenant_id` can be retrieved after the Function App has been created. More details are available below.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned. Required if `type` is `UserAssigned`.

---

An `auth_settings` block supports the following:

* `enabled` - (Required) Is Authentication enabled?

* `active_directory` - (Optional) An `active_directory` block as defined below.

* `additional_login_params` - (Optional) login parameters to send to the OpenID Connect authorization endpoint when a user logs in. Each parameter must be in the form "key=value".

* `allowed_external_redirect_urls` - (Optional) External URLs that can be redirected to as part of logging in or logging out of the app.

* `default_provider` - (Optional) The default provider to use when multiple providers have been set up. Possible values are `AzureActiveDirectory`, `Facebook`, `Google`, `MicrosoftAccount` and `Twitter`.

~> **Note:** When using multiple providers, the default provider must be set for settings like `unauthenticated_client_action` to work.

* `facebook` - (Optional) A `facebook` block as defined below.

* `google` - (Optional) A `google` block as defined below.

* `issuer` - (Optional) Issuer URI. When using Azure Active Directory, this value is the URI of the directory tenant, e.g. <https://sts.windows.net/{tenant-guid}/>.

* `microsoft` - (Optional) A `microsoft` block as defined below.

* `runtime_version` - (Optional) The runtime version of the Authentication/Authorization module.

* `token_refresh_extension_hours` - (Optional) The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to `72`.

* `token_store_enabled` - (Optional) If enabled the module will durably store platform-specific security tokens that are obtained during login flows. Defaults to `false`.

* `twitter` - (Optional) A `twitter` block as defined below.

* `unauthenticated_client_action` - (Optional) The action to take when an unauthenticated client attempts to access the app. Possible values are `AllowAnonymous` and `RedirectToLoginPage`.

---

An `active_directory` block supports the following:

* `client_id` - (Required) The Client ID of this relying party application. Enables OpenIDConnection authentication with Azure Active Directory.

* `client_secret` - (Optional) The Client Secret of this relying party application. If no secret is provided, implicit flow will be used.

* `allowed_audiences` - (Optional) Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

---

A `facebook` block supports the following:

* `app_id` - (Required) The App ID of the Facebook app used for login

* `app_secret` - (Required) The App Secret of the Facebook app used for Facebook login.

* `oauth_scopes` - (Optional) The OAuth 2.0 scopes that will be requested as part of Facebook login authentication. <https://developers.facebook.com/docs/facebook-login>

---

A `google` block supports the following:

* `client_id` - (Required) The OpenID Connect Client ID for the Google web application.

* `client_secret` - (Required) The client secret associated with the Google web application.

* `oauth_scopes` - (Optional) The OAuth 2.0 scopes that will be requested as part of Google Sign-In authentication. <https://developers.google.com/identity/sign-in/web/>

---

A `microsoft` block supports the following:

* `client_id` - (Required) The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret` - (Required) The OAuth 2.0 client secret that was created for the app used for authentication.

* `oauth_scopes` - (Optional) The OAuth 2.0 scopes that will be requested as part of Microsoft Account authentication. <https://msdn.microsoft.com/en-us/library/dn631845.aspx>

---

A `twitter` block supports the following:

* `consumer_key` - (Required) The OAuth 1.0a consumer key of the Twitter application used for sign-in.

* `consumer_secret` - (Required) The OAuth 1.0a consumer secret of the Twitter application used for sign-in.

---

A `ip_restriction` block supports the following:

* `ip_address` - (Optional) The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

-> **Note:** One of either `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified

* `name` - (Optional) The name for this IP Restriction.

* `priority` - (Optional) The priority for this IP Restriction. Restrictions are enforced in priority order. By default, priority is set to 65000 if not specified.

* `action` - (Optional) Does this restriction `Allow` or `Deny` access for this IP range. Defaults to `Allow`. 

* `headers` - (Optional) The `headers` block for this specific `ip_restriction` as defined below.

---

A `scm_ip_restriction` block supports the following:

* `ip_address` - (Optional) The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

-> **Note:** One of either `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified

* `name` - (Optional) The name for this IP Restriction.

* `priority` - (Optional) The priority for this IP Restriction. Restrictions are enforced in priority order. By default, priority is set to 65000 if not specified.

* `action` - (Optional) Allow or Deny access for this IP range. Defaults to `Allow`.

* `headers` - (Optional) The `headers` block for this specific `scm_ip_restriction` as defined below.

---

A `headers` block supports the following:

* `x_azure_fdid` - (Optional) A list of allowed Azure FrontDoor IDs in UUID notation with a maximum of 8.

* `x_fd_health_probe` - (Optional) A list to allow the Azure FrontDoor health probe header. Only allowed value is "1".

* `x_forwarded_for` - (Optional) A list of allowed 'X-Forwarded-For' IPs in CIDR notation with a maximum of 8

* `x_forwarded_host` - (Optional) A list of allowed 'X-Forwarded-Host' domains with a maximum of 8.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Function App Slot

* `default_hostname` - The default hostname associated with the Function App - such as `mysite.azurewebsites.net`

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12`

* `possible_outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12,52.143.43.17` - not all of which are necessarily in use. Superset of `outbound_ip_addresses`.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Function App Slot.

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this Function App Slot.

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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Function App Deployment Slot.
* `read` - (Defaults to 5 minutes) Used when retrieving the Function App Deployment Slot.
* `update` - (Defaults to 30 minutes) Used when updating the Function App Deployment Slot.
* `delete` - (Defaults to 30 minutes) Used when deleting the Function App Deployment Slot.

## Import

Function Apps Deployment Slots can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_function_app_slot.functionapp1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/functionapp1/slots/staging
```
